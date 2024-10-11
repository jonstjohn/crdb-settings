package metrics

import (
	"fmt"
	"github.com/jonstjohn/crdb-settings/pkg/crdbcluster"
	"github.com/jonstjohn/crdb-settings/pkg/releases"
	"github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
)

type Manager struct {
	Db *Db
}

func NewManager(url string) (*Manager, error) {

	db, err := NewDbDatasource(url)
	if err != nil {
		return nil, err
	}
	return &Manager{Db: db}, err
}

func (m *Manager) InitializeDatabase() error {
	err := m.Db.createRawTableIfNotExists()
	if err != nil {
		return err
	}
	err = m.Db.createSaveRunsTableIfNotExists()
	return err
}

func (m *Manager) GetMetricsForRelease(releaseName string) ([]Metric, error) {
	rows, err := m.Db.SelectRaw(releaseName)
	if err != nil {
		return nil, err
	}

	ms := make([]Metric, 0)
	for _, row := range rows {
		ms = append(ms, Metric{Name: row.Metric, Help: row.Help, Type: Type(row.Type)})
	}
	return ms, err
}

func (m *Manager) SaveClusterSettingsForRelease(releaseName string) error {
	rs, err := m.getReleasesNames(releaseName)
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("Found %d releases that are candidate for updating", len(rs)))

	cm := crdbcluster.NewManager()

	// Iterate over releases
	for _, r := range rs {
		runs, err := m.Db.SelectSaveRuns(r)
		if err != nil {
			return err
		}
		if len(runs) > 0 {
			logrus.Info(fmt.Sprintf("Save run already exists for '%s', skipping", r))
			continue
		}

		// Start test server
		err := cm.StartTestCluster(r)
		if err != nil {
			return err
		}

		// Next, get the metrics for the release TODO

		// Save it

		// Cleanup test server
		err = cm.CleanupTestCluster()
		if err != nil {
			return err
		}
	}
}

func (m *Manager) getReleasesNames(release string) ([]string, error) {
	if release == "all" || strings.HasPrefix(release, "recent-") {
		cnt := math.MaxInt
		var err error
		if strings.HasPrefix(release, "recent-") {
			cntStr := strings.Replace(release, "recent-", "", 1)
			cnt, err = strconv.Atoi(cntStr)
			if err != nil {
				return nil, err
			}
		}
		rm, err := releases.NewReleasesManager(m.Db.Url)
		if err != nil {
			return nil, err
		}
		return rm.GetRecentReleaseNames(cnt)
	} else {
		return []string{release}, nil
	}
}
