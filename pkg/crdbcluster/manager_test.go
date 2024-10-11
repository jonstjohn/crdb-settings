package crdbcluster

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManager_GetMetricsEndpoint(t *testing.T) {
	m := NewManager()
	err := m.StartTestCluster("v23.2.10")
	defer func(m *Manager) {
		err := m.CleanupTestCluster()
		if err != nil {
			assert.NoError(t, err)
		}
	}(m)

	assert.NoError(t, err)

	metricsUrl, err := m.GetMetricsEndpoint()
	assert.NoError(t, err)

	assert.NotEmpty(t, metricsUrl)
}

func TestManager_GetMetricsEndpointOutput(t *testing.T) {
	m := NewManager()
	m.StartTestCluster("v23.2.10")
	defer func(m *Manager) {
		err := m.CleanupTestCluster()
		if err != nil {
			assert.NoError(t, err)
		}
	}(m)

	output, err := m.GetMetricsEndpointOutput()
	assert.NoError(t, err)

	assert.Contains(t, string(output), "HELP")
}
