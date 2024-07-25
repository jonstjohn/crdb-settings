package settings

import (
	"slices"
	"strings"
)

type ChangedSettings []ChangedSetting

type ChangedSetting struct {
	Before ReleaseSetting `json:"before"`
	After  ReleaseSetting `json:"after"`
}

type ComparedReleaseSettings struct {
	Added   ReleaseSettings `json:"added"`
	Removed ReleaseSettings `json:"removed"`
	Changed ChangedSettings `json:"changed"`
}

func CompareReleaseSettings(rs1 ReleaseSettings, rs2 ReleaseSettings) ComparedReleaseSettings {
	rs1indexed := make(map[string]ReleaseSetting)
	for _, rs := range rs1 {
		rs1indexed[rs.Variable] = rs
	}
	rs2indexed := make(map[string]ReleaseSetting)
	for _, rs := range rs2 {
		rs2indexed[rs.Variable] = rs
	}

	added := ReleaseSettings{}
	removed := ReleaseSettings{}
	changed := ChangedSettings{}

	for _, r1 := range rs1 {
		if slices.Contains(IgnoredSettings, r1.Variable) {
			continue
		}
		if _, ok := rs2indexed[r1.Variable]; !ok { // exists in r1 but not r2
			removed = append(removed, r1)
			continue
		}
		if r2, ok := rs2indexed[r1.Variable]; ok {
			if r1.Value != r2.Value || strings.TrimSpace(r1.Description) != strings.TrimSpace(r2.Description) {
				changed = append(changed, ChangedSetting{Before: r1, After: r2})
				continue
			}
		}
	}

	for _, r2 := range rs2 {
		if slices.Contains(IgnoredSettings, r2.Variable) {
			continue
		}
		if _, ok := rs1indexed[r2.Variable]; !ok { // exists in r2 but not r1
			added = append(added, r2)
			continue
		}
	}

	return ComparedReleaseSettings{
		Removed: removed,
		Added:   added,
		Changed: changed,
	}
}
