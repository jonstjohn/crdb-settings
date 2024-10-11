package metrics

import (
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManager_GetMetricsFromClusterForRelease(t *testing.T) {
	m, err := NewManager("")
	assert.NoError(t, err)

	metrics, err := m.GetMetricsFromClusterForRelease("v23.2.10")
	assert.NoError(t, err)

	assert.Equal(t, 1578, len(metrics))
	assert.Equal(t, "abortspanbytes", metrics[0].Name)
	assert.Equal(t, Type("gauge"), metrics[0].Type)
	assert.Equal(t, "Number of bytes in the abort span", metrics[0].Help)

}

func TestManager_SaveMetricsForRelease(t *testing.T) {
	ts, err := testserver.NewTestServer(
		testserver.CustomVersionOpt("v23.2.10"))
	assert.NoError(t, err)
	err = ts.Start()
	assert.NoError(t, err)
	url := ts.PGURL().String()

	m, err := NewManager(url)
	assert.NoError(t, err)
	assert.NoError(t, m.InitializeDatabase())
	err = m.SaveMetricsForRelease("v23.2.10")
	assert.NoError(t, err)

	metrics, err := m.GetMetrics("v23.2.10")
	assert.NoError(t, err)
	assert.Equal(t, "abortspanbytes", metrics[0].Name)
	assert.Equal(t, Type("gauge"), metrics[0].Type)
	assert.Equal(t, "Number of bytes in the abort span", metrics[0].Help)

}
