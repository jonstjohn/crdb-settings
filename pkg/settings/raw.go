package settings

import (
	"github.com/jonstjohn/crdb-settings/pkg/releases"
	"slices"
	"time"
)

type RawSettings []RawSetting

type RawSetting struct {
	ReleaseName  string
	Cpu          int
	MemoryBytes  int64
	Variable     string
	Value        string
	Type         string
	Public       bool
	Description  string
	DefaultValue string
	Origin       string
	Key          string
	Updated      time.Time
}

var ignoreList = []string{"cluster.secret", "version"}

// TODO Okay. I know this isn't great but trying to be pragmatic here for now
type RawSettingsWithReleases []RawSettingWithRelease
type RawSettingWithRelease struct {
	RawSetting *RawSetting
	Release    *releases.Release
}

type Meta struct {
	mostRecent         RawSetting
	firstReleases      []string
	lastReleases       []string
	hostDependent      bool
	valueChanges       []Change
	descriptionChanges []Change
}

func NewRawSetting(releaseName string, cpu int, memoryBytes int64, cs ClusterSetting) *RawSetting {
	return &RawSetting{
		ReleaseName:  releaseName,
		Cpu:          cpu,
		MemoryBytes:  memoryBytes,
		Variable:     cs.Variable,
		Value:        cs.Value,
		Type:         cs.Type,
		Public:       cs.Public,
		Description:  cs.Description,
		DefaultValue: cs.DefaultValue,
		Origin:       cs.Origin,
		Key:          cs.Key,
	}
}

func (r *RawSetting) Compare(r2 *RawSetting) int {
	if r.Variable == r2.Variable {
		if r.ReleaseName == r2.ReleaseName {
			if r.Cpu == r2.Cpu {
				return 0
			} else if r.Cpu < r2.Cpu {
				return -1
			} else {
				return 1
			}
		} else if r.ReleaseName < r2.ReleaseName {
			return -1
		} else {
			return 1
		}
	} else if r.Variable < r2.Variable {
		return -1
	} else {
		return 1
	}
}

func (rs *RawSettings) Sort() {
	slices.SortFunc(*rs, func(a, b RawSetting) int {
		return a.Compare(&b)
	})
}

func (rs *RawSettings) SortedVariables() []string {
	vs := make([]string, 0)
	m := make(map[string]bool)
	for _, r := range *rs {
		if !m[r.Variable] {
			m[r.Variable] = true
			vs = append(vs, r.Variable)
		}
	}
	slices.Sort(vs)
	return vs
}

// LatestForVariable returns the latest raw setting for a specific variable
func (rs *RawSettings) LatestForVariable(v string) *RawSetting {
	rs.Sort() // TODO lots of extra sorting
	lastR := RawSetting{}
	for i, r := range *rs {
		if i > 0 && r.Variable != lastR.Variable {
			return &lastR
		}
		lastR = r
	}
	return nil
}

func (rs *RawSettings) ForVariableOnly(v string) *RawSettings {
	rsv := RawSettings{}
	for _, r := range *rs {
		if r.Variable == v {
			rsv = append(rsv, r)
		}
	}
	return &rsv
}

func (rs *RawSettings) ReleaseNames() []string {
	m := make(map[string]bool)
	names := make([]string, 0)
	for _, r := range *rs {
		if !m[r.ReleaseName] {
			m[r.ReleaseName] = true
			names = append(names, r.ReleaseName)
		}
	}
	return names
}

// MetaForVariable takes a variable name and a set of releases then creates some meta data about the variable
// across the raw settings which can be used for setting summaries
func (rs *RawSettings) MetaForVariable(v string, rels releases.Releases) *Meta {

	// Extract and sort rswr settings for single variable
	rsv := rs.ForVariableOnly(v)
	//rsv.Sort()

	rswrs := NewRawSettingsWithReleases(*rsv, rels)
	rswrs.SortByRelease()

	names := rsv.ReleaseNames()

	releasesForVariable := make(releases.Releases, 0)
	for _, r := range rels {
		if !slices.Contains(names, r.Name) {
			continue
		}
		releasesForVariable = append(releasesForVariable, r)
	}
	firstReleasesPerMajor := releasesForVariable.FirstReleasePerMajorVersion()
	lastReleasesPerMajor := releasesForVariable.LastReleasePerMajorVersion()

	firstReleases := make([]string, len(firstReleasesPerMajor))
	lastReleases := make([]string, len(lastReleasesPerMajor))

	for i, fr := range firstReleasesPerMajor {
		firstReleases[i] = fr.Name
	}

	for i, lr := range lastReleasesPerMajor {
		lastReleases[i] = lr.Name
	}

	// Iterate over rswr settings to see if setting has changed across versions
	type releaseValue struct {
		Release string
		Value   string
	}
	type releaseDescription struct {
		Release     string
		Description string
	}
	currentValueForCpu := make(map[int]releaseValue)             // map of CPU to value
	currentDescriptionForCpu := make(map[int]releaseDescription) // map of CPU to description
	valueChanges := make([]Change, 0)
	descriptionChanges := make([]Change, 0)
	valuesForRelease := make(map[string][]string)
	hostDependent := false
	for _, rswr := range rswrs {
		// Initialize current value if it hasn't been set
		if _, ok := currentValueForCpu[rswr.RawSetting.Cpu]; !ok {
			currentValueForCpu[rswr.RawSetting.Cpu] = releaseValue{
				Release: rswr.RawSetting.ReleaseName, Value: rswr.RawSetting.Value}
			continue
		}

		// Initialize current description if it hasn't been set
		if _, ok := currentDescriptionForCpu[rswr.RawSetting.Cpu]; !ok {
			currentDescriptionForCpu[rswr.RawSetting.Cpu] = releaseDescription{Release: rswr.RawSetting.ReleaseName, Description: rswr.RawSetting.Description}
			continue
		}

		// If the value has changed for the same CPU, record it as a value change
		if rswr.RawSetting.Value != currentValueForCpu[rswr.RawSetting.Cpu].Value {
			if !slices.ContainsFunc(valueChanges, func(c Change) bool {
				return c.Release == rswr.RawSetting.ReleaseName
			}) {
				valueChanges = append(valueChanges,
					Change{
						Release: rswr.RawSetting.ReleaseName,
						From:    currentValueForCpu[rswr.RawSetting.Cpu].Value,
						To:      rswr.RawSetting.Value,
					})
			}
			currentValueForCpu[rswr.RawSetting.Cpu] = releaseValue{Release: rswr.RawSetting.ReleaseName, Value: rswr.RawSetting.Value}
		}

		// If the description has changed for the same CPU, record it as a value change
		if rswr.RawSetting.Description != currentDescriptionForCpu[rswr.RawSetting.Cpu].Description {
			if !slices.ContainsFunc(descriptionChanges, func(c Change) bool {
				return c.Release == rswr.RawSetting.ReleaseName
			}) {
				descriptionChanges = append(descriptionChanges,
					Change{
						Release: rswr.RawSetting.ReleaseName,
						From:    currentDescriptionForCpu[rswr.RawSetting.Cpu].Description,
						To:      rswr.RawSetting.Description,
					})
			}
			currentDescriptionForCpu[rswr.RawSetting.Cpu] = releaseDescription{Release: rswr.RawSetting.ReleaseName, Description: rswr.RawSetting.Description}
		}

		// If the value for the same release and a different CPU is not the same, mark it as host dependent
		if _, ok := valuesForRelease[rswr.RawSetting.ReleaseName]; !ok {
			valuesForRelease[rswr.RawSetting.ReleaseName] = []string{rswr.RawSetting.Value}
		} else { // otherwise, check to see if there is a different value for this release
			if !slices.Contains(valuesForRelease[rswr.RawSetting.ReleaseName], rswr.RawSetting.Value) {
				hostDependent = true
			}
			valuesForRelease[rswr.RawSetting.ReleaseName] = append(valuesForRelease[rswr.RawSetting.ReleaseName], rswr.RawSetting.Value)
		}

	}

	return &Meta{
		mostRecent:         *rswrs[len(rswrs)-1].RawSetting,
		firstReleases:      firstReleases,
		lastReleases:       lastReleases,
		valueChanges:       valueChanges,
		descriptionChanges: descriptionChanges,
		hostDependent:      hostDependent,
	}

}

func NewRawSettingsWithReleases(rawSettings RawSettings, rels releases.Releases) RawSettingsWithReleases {
	m := make(map[string]*releases.Release)
	for _, r := range rels {
		m[r.Name] = &r
	}

	rswr := make([]RawSettingWithRelease, 0)
	for _, setting := range rawSettings {
		rswr = append(rswr, RawSettingWithRelease{
			RawSetting: &setting,
			Release:    m[setting.ReleaseName],
		})
	}
	return rswr
}

func (rss RawSettingsWithReleases) SortByRelease() {
	slices.SortFunc(rss, func(a, b RawSettingWithRelease) int {
		return a.Release.CompareVersion(b.Release)
	})
}
