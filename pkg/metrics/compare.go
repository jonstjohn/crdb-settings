package metrics

type ReleaseMetric struct {
	Release string `json:"release"`
	Metric  string `json:"metric"`
	Help    string `json:"help"`
	Type    Type   `json:"type"`
}

type ChangedMetrics []ChangedMetric

type ChangedMetric struct {
	Before ReleaseMetric `json:"before"`
	After  ReleaseMetric `json:"after"`
}

type ComparedReleaseMetrics struct {
	Added   Metrics        `json:"added"`
	Removed Metrics        `json:"removed"`
	Changed ChangedMetrics `json:"changed"`
}

func CompareReleaseMetrics(r1 string, r1metrics Metrics, r2 string, r2metrics Metrics) ComparedReleaseMetrics {

	rs1indexed := make(map[string]Metric)
	for _, rs := range r1metrics {
		rs1indexed[rs.Name] = rs
	}
	rs2indexed := make(map[string]Metric)
	for _, rs := range r2metrics {
		rs2indexed[rs.Name] = rs
	}

	added := Metrics{}
	removed := Metrics{}
	var changed ChangedMetrics

	for _, r1m := range r1metrics {
		if _, ok := rs2indexed[r1m.Name]; !ok { // exists in r1 but not r2
			removed = append(removed, r1m)
			continue
		}

		if r2m, ok := rs2indexed[r1m.Name]; !ok {
			if r1m.Name != r2m.Name || r1m.Help != r2m.Help {
				changed = append(changed, ChangedMetric{
					Before: ReleaseMetric{
						Release: r1,
						Metric:  r1m.Name,
						Help:    r1m.Help,
						Type:    r1m.Type,
					},
					After: ReleaseMetric{
						Release: r2,
						Metric:  r2m.Name,
						Help:    r2m.Help,
						Type:    r2m.Type,
					},
				})
			}
		}
	}

	for _, r2m := range r2metrics {
		if _, ok := rs1indexed[r2m.Name]; !ok {
			added = append(added, r2m)
		}
	}

	return ComparedReleaseMetrics{
		Removed: removed,
		Added:   added,
		Changed: changed,
	}
}
