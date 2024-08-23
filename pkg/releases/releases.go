package releases

import (
	"slices"
	"time"
)

// The releases package provides functionality around CockroachDB releases, including listing and saving to
// a remote cluster

type Releases []Release

type SortBy int

const (
	SortByVersion         SortBy = iota
	SortByReleaseDate     SortBy = iota
	SortByVersionReversed SortBy = iota
)

type Release struct {
	Name          string    `json:"release_name"`
	Withdrawn     bool      `json:"withdrawn"`
	CloudOnly     bool      `json:"cloud_only"`
	ReleaseType   string    `json:"release_type"` // Production or Testing
	ReleaseDate   time.Time `json:"release_date"`
	MajorVersion  string    `json:"major_version"` // e.g., v1.0, v23.1
	Major         int       `json:"major"`
	Minor         int       `json:"minor"`
	Patch         int       `json:"path"`
	BetaRc        string    `json:"beta_rc"`
	BetaRcVersion int       `json:"beta_rc_version"`
}

func (r *Release) CompareDates(r2 *Release) int {
	if r.ReleaseDate == r2.ReleaseDate {
		return 0
	} else if r.ReleaseDate.Before(r2.ReleaseDate) {
		return -1
	}
	return 1
}

func (r *Release) CompareVersion(r2 *Release) int {
	// Check to see if r is < r2 and default to false
	if r.Major < r2.Major { // Major is less
		return -1
	} else if r.Major == r2.Major { // Majors are equal
		if r.Minor < r2.Minor { // Minor is less
			return -1
		} else if r.Minor == r2.Minor { // Minors are equal
			if r.Patch < r2.Patch { // Patch is less
				return -1
			} else if r.Patch == r2.Patch { // Patches are equal
				if r.BetaRc == "alpha" && (r2.BetaRc == "beta" || r2.BetaRc == "rc" || r2.BetaRc == "") { // Alpha v
					return -1
				} else if r.BetaRc == "beta" && (r2.BetaRc == "rc" || r2.BetaRc == "") { // Beta vs RC or prod
					return -1
				} else if r.BetaRc == "rc" && r2.BetaRc == "" { // RC vs prod
					return -1
				} else if r.BetaRc == r2.BetaRc { // Same alpha, beta or rc, look at version
					if r.BetaRcVersion < r2.BetaRcVersion {
						return -1
					} else if r.BetaRcVersion == r2.BetaRcVersion {
						return 0
					} else {
						return 1
					}
				} else {
					return 1
				}
			} else {
				return 1
			}
		} else {
			return 1
		}
	} else {
		return 1
	}
	return 1
}

func (rs *Releases) SortBy(sort SortBy) {
	switch sort {
	case SortByVersion:
		slices.SortFunc(*rs, func(a, b Release) int {
			return a.CompareVersion(&b)
		})
	case SortByVersionReversed:
		slices.SortFunc(*rs, func(a, b Release) int {
			return 0 - a.CompareVersion(&b)
		})
	case SortByReleaseDate:
		slices.SortFunc(*rs, func(a, b Release) int {
			return a.CompareDates(&b)
		})
	}
}

// MajorVersions returns the major versions for all of the releases sorted by version
func (rs *Releases) MajorVersions() []string {
	rs.SortBy(SortByVersion)
	mvs := make([]string, 0)
	for _, r := range *rs {
		if !slices.Contains(mvs, r.MajorVersion) {
			mvs = append(mvs, r.MajorVersion)
		}
	}
	return mvs
}

func (rs *Releases) FirstTestingReleaseForMajorVersion(mv string) *Release {
	for _, r := range *rs {
		if r.MajorVersion == mv && r.ReleaseType == "Testing" {
			return &r
		}
	}
	return nil
}

func (rs *Releases) FirstProductionReleaseForMajorVersion(mv string) *Release {
	for _, r := range *rs {
		if r.MajorVersion == mv && r.ReleaseType == "Production" {
			return &r
		}
	}
	return nil
}

func (rs *Releases) LatestReleaseForMajorVersion(mv string) *Release {
	lastR := Release{}
	for _, r := range *rs {
		if lastR.MajorVersion == mv && r.MajorVersion != mv {
			return &r
		}
		lastR = r
	}
	if lastR.MajorVersion == mv {
		return &lastR
	}
	return nil
}

func (rs *Releases) GetReleaseForName(name string) *Release {
	for _, r := range *rs {
		if r.Name == name {
			return &r
		}
	}

	return nil
}

func (rs *Releases) FilterForNames(names []string) (ret Releases) {
	for _, r := range *rs {
		if slices.Contains(names, r.Name) {
			ret = append(ret, r)
		}
	}
	return
}

func (rs *Releases) MostRecent() Release {
	rs.SortBy(SortByVersion)
	return (*rs)[len(*rs)-1]
}

// FirstReleasePerMajorVersion returns a list of releases that were the first for the set of releases
func (rs *Releases) FirstReleasePerMajorVersion() []Release {
	rs.SortBy(SortByVersion)
	rsmv := make([]Release, 0)
	m := make(map[string]bool)
	for _, r := range *rs {
		if _, ok := m[r.MajorVersion]; !ok {
			m[r.MajorVersion] = true
			rsmv = append(rsmv, r)
		}
	}

	return rsmv
}

// LastReleasePerMajorVersion returns a list of releases that were the last for the set of releases
func (rs *Releases) LastReleasePerMajorVersion() Releases {
	rs.SortBy(SortByVersionReversed)
	rsmv := make(Releases, 0)
	m := make(map[string]bool)
	for _, r := range *rs {
		if _, ok := m[r.MajorVersion]; !ok {
			m[r.MajorVersion] = true
			rsmv = append(rsmv, r)
		}
	}

	rsmv.SortBy(SortByVersion)
	return rsmv
}
