package settings

type Provider interface {
	GetRawSettings() (RawSettings, error)
	GetSettingsSummary() (Summaries, error)
}

type Persister interface {
	SaveRawSettings(RawSettings) error
	SaveSettingsSummaries(Summaries) error
	SaveRun(string, int, int64) error
}
