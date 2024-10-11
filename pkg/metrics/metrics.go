package metrics

type Metrics []Metric

type Metric struct {
	Name string
	Help string
	Type Type
}
