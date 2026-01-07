package crdbcluster

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/jonstjohn/crdb-settings/pkg/dbpgx"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	TestServer *testserver.TestServer
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) StartTestCluster(releaseName string) error {
	t, err := testserver.NewTestServer(
		testserver.CustomVersionOpt(releaseName))
	if err != nil {
		return err
	}

	m.TestServer = &t
	return nil
}

func (m *Manager) CleanupTestCluster() error {

	(*m.TestServer).Stop()

	// Delete files
	files, err := filepath.Glob(filepath.Join(os.TempDir(), "cockroach-v*"))
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}

	m.TestServer = nil

	return nil
}

func (m *Manager) GetPGUrl() (*url.URL, error) {
	if err := m.errorTestServerNotRunning(); err != nil {
		return nil, err
	}
	return (*m.TestServer).PGURL(), nil
}

func (m *Manager) GetDbConsoleURL() (*url.URL, error) {
	if err := m.errorTestServerNotRunning(); err != nil {
		return nil, err
	}

	pgurl, err := m.GetPGUrl()
	if err != nil {
		return nil, err
	}
	pool, err := dbpgx.NewPoolFromUrl(pgurl.String())

	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return nil, err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "SET allow_unsafe_internals = 'on'")
	if err != nil {
		logrus.Warnf("could not enable unsafe internals: %v", err)
	}

	sql := `
SELECT value
FROM   crdb_internal.node_runtime_info
WHERE component = 'UI' AND field = 'URL'
LIMIT 1
`

	row := conn.QueryRow(context.Background(), sql)

	var urlStr string
	err = row.Scan(&urlStr)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return url, nil

}

func (m *Manager) GetMetricsEndpointOutput() ([]byte, error) {
	ep, err := m.GetMetricsEndpoint()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(ep.String())
	if err != nil {
		return nil, fmt.Errorf("could not download metrics data: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (m *Manager) GetMetricsEndpoint() (*url.URL, error) {
	if err := m.errorTestServerNotRunning(); err != nil {
		return nil, err
	}
	consoleUrl, err := m.GetDbConsoleURL()
	if err != nil {
		return nil, err
	}
	metricsUrl := consoleUrl
	metricsUrl.Path = path.Join(consoleUrl.Path, "_status", "vars")
	return metricsUrl, nil
}

func (m *Manager) errorTestServerNotRunning() error {
	if m.TestServer == nil {
		return fmt.Errorf("test server is not running")
	}
	return nil
}
