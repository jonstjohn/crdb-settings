package settings

import "github.com/jonstjohn/crdb-settings/pkg/releases"

// Summary - create a summary of all the raw settings. Here is what we're interested in:
// - variable name
// - default values over releases
// -

// Iterate over every setting
// Track the first version that it appeared in, if there is a difference between value with different CPU,
// and when the value changes across versions, and if it doesn't exist in the current version
/*
select settings_raw.release_name, settings_raw.cpu,
	settings_raw.memory_bytes, settings_raw.variable, settings_raw.value
from
	settings_raw INNER JOIN releases ON settings_raw.release_name = releases.name
order
	by settings_raw.variable, releases.major, releases.minor, releases.patch,
	releases.beta_rc, releases.beta_rc_version, settings_raw.cpu,
	settings_raw.memory_bytes
limit 100;
*/

type Summarizer struct {
	RawSettings RawSettings
	Releases    releases.Releases
}

type Change struct {
	Release string
	From    string
	To      string
}

type Summaries []Summary

type Summary struct {
	Variable           string
	Value              string
	Type               string
	Public             bool
	Description        string
	DefaultValue       string
	Origin             string
	Key                string
	FirstReleases      []string
	LastReleases       []string
	HostDependent      bool
	ValueChanges       []Change
	DescriptionChanges []Change
}

func NewSummarizer(rawSettings RawSettings, rels releases.Releases) *Summarizer {
	return &Summarizer{RawSettings: rawSettings, Releases: rels}
}

// SummarizeAndSave takes raw settings and releaeses, creating
func (sum *Summarizer) Summarize() ([]Summary, error) {

	sum.RawSettings.Sort()
	variables := sum.RawSettings.SortedVariables()
	summaries := make([]Summary, 0)
	for _, v := range variables {
		meta := sum.RawSettings.MetaForVariable(v, sum.Releases)
		s := Summary{
			Variable:           meta.mostRecent.Variable,
			Value:              meta.mostRecent.Value,
			Type:               meta.mostRecent.Type,
			Public:             meta.mostRecent.Public,
			Description:        meta.mostRecent.Description,
			DefaultValue:       meta.mostRecent.DefaultValue,
			Origin:             meta.mostRecent.Origin,
			Key:                meta.mostRecent.Key,
			FirstReleases:      meta.firstReleases,
			LastReleases:       meta.lastReleases,
			HostDependent:      meta.hostDependent,
			ValueChanges:       meta.valueChanges,
			DescriptionChanges: meta.descriptionChanges,
		}
		summaries = append(summaries, s)
	}

	return summaries, nil
}
