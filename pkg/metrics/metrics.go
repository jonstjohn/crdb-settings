package metrics

type Metrics []Metric

type Metric struct {
	Name string `json:"name"`
	Help string `json:"help"`
	Type Type   `json:"type"`
}
