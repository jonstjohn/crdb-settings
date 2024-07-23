package settings

type SettingChangeType int

const (
	FirstSeen SettingChangeType = iota
	LastSeen
	Changed
)

type SettingHistory struct {
	Change SettingHistoryChange
}

type SettingHistoryChange struct {
	FromRelease string
	ToRelease   string
	Before      *ReleaseSetting
	After       *ReleaseSetting
}

func GenerateSettingHistory(settings ReleaseSettings) (SettingHistory, error) {
	return SettingHistory{}, nil
}
