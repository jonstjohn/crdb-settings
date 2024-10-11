package metrics

import (
	"bufio"
	"fmt"
	"strings"
)

type Type string

const (
	Counter   Type = "counter"
	Gauge          = "gauge"
	Histogram      = "histogram"
)

type LineType int

const (
	Help LineType = iota
	MetricType
	Data
)

func FromText(text string) Metrics {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var section []string
	var metrics Metrics

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# HELP") {
			if section != nil {
				metrics = append(metrics, parseSection(section))
			}
			section = nil
			section = append(section, line)
		} else {
			section = append(section, line)
		}
	}

	return metrics
}

func parseSection(section []string) Metric {
	m := Metric{}
	for _, line := range section {
		if strings.HasPrefix(line, "# HELP") {
			name, help, _ := strings.Cut(line[len("# HELP "):], " ")
			m.Name = name
			m.Help = help
		} else if strings.HasPrefix(line, "# TYPE") {
			_, typ, _ := strings.Cut(line[len(fmt.Sprintf("# TYPE %s", m.Name)):], " ")
			m.Type = Type(typ)
		}
	}
	return m
}
