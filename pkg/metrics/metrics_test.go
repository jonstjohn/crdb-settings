package metrics

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestMetricsFromOutput(t *testing.T) {
	file, err := os.Open("testdata/23_2_10_metrics.txt")
	assert.NoError(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var s []string
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}

	metrics := FromText(strings.Join(s, "\n"))
	assert.Equal(t, 1597, len(metrics))
	assert.Equal(t, "kv_rangefeed_budget_allocation_failed", metrics[0].Name)
	assert.Equal(t, "Number of times RangeFeed failed because memory budget was exceeded", metrics[0].Help)
	assert.Equal(t, Counter, metrics[0].Type)

}
