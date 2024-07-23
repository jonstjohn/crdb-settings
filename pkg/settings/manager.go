package settings

import (
	"fmt"
	"github.com/jonstjohn/crdb-settings/pkg/host"
	"github.com/jonstjohn/crdb-settings/pkg/releases"
	"github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
)

type Manager struct {
	Db *Db
}

func NewSettingsManager(url string) (*Manager, error) {

	db, err := NewDbDatasource(url)
	if err != nil {
		return nil, err
	}
	return &Manager{Db: db}, err
}

func (sm *Manager) GetSettingsForRelease(version string) (ReleaseSettings, error) {
	raws, err := sm.Db.GetRawSettingsForVersion(version)
	s := make(ReleaseSettings, len(raws))
	if err != nil {
		return s, err // TODO
	}
	for i, raw := range raws {
		s[i] = ReleaseSetting{
			ReleaseName: raw.ReleaseName,
			Variable:    raw.Variable,
			Value:       raw.Value,
			Type:        raw.Type,
			Public:      raw.Public,
			Description: raw.Description,
		}
	}
	return s, nil
}

// SaveClusterSettingsForVersion saves all the cluster settings for a specific CRDB version, but only
// if the combination of release, cpu and memory has not been previously run - otherwise it bails early.
func (sm *Manager) SaveClusterSettingsForVersion(release string, url string) error {

	// Get host memory and CPU
	cpu := host.GetCpu()
	memoryBytes, err := host.GetMemory()
	if err != nil {
		return err
	}

	rs, err := sm.getReleasesNames(release)
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("Found %d releases that are candidate for updating", len(rs)))

	// Iterate over releases
	for _, r := range rs {

		// Check to see if save run already exists, if it does, bail early - we've already captured the settings
		exists, err := sm.Db.SaveRunExists(r, cpu, memoryBytes) // TODO
		if err != nil {
			return err
		}
		if exists {
			logrus.Info(fmt.Sprintf("Save run already exists for '%s' with cpu/memory %d/%d", r, cpu, memoryBytes))
			continue
		}

		// Get the cluster settings for this release
		settings, err := ClusterSettingsFromRelease(r)
		if err != nil {
			return err
		}
		rawSettings := make([]RawSetting, len(settings))

		// Convert the cluster settings into raw settings to be saved
		for i, s := range settings {
			rawSettings[i] = *NewRawSetting(r, cpu, memoryBytes, s)
		}

		err = sm.Db.SaveRawSettings(rawSettings) // TODO
		if err != nil {
			return err
		}

		// Save the save run so we don't have to re-run later
		err = sm.Db.SaveRun(r, cpu, memoryBytes) // TODO
		if err != nil {
			return err
		}
	}

	return nil

}

func (sm *Manager) getReleasesNames(release string) ([]string, error) {
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
		rdb := releases.NewDbDatasource(sm.Db.Pool) // TODO
		return rdb.GetRecentReleaseNames(cnt)
	} else {
		return []string{release}, nil
	}
}

func (sm *Manager) CompareSettingsForReleases(r1 string, r2 string) (ComparedReleaseSettings, error) {
	rs1, err := sm.GetSettingsForRelease(r1)
	if err != nil {
		return ComparedReleaseSettings{}, err
	}

	rs2, err := sm.GetSettingsForRelease(r2)
	if err != nil {
		return ComparedReleaseSettings{}, err
	}

	return CompareReleaseSettings(rs1, rs2), nil

}
