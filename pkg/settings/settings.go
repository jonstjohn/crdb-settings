package settings

import (
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/jonstjohn/crdb-settings/pkg/dbpgx"
	"github.com/jonstjohn/crdb-settings/pkg/releases"
	"os"
	"path/filepath"
)

type ReleaseSettings []ReleaseSetting

type ReleaseSetting struct {
	ReleaseName string `json:"release_name""`
	Variable    string `json:"variable"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Public      bool   `json:"public"`
	Description string `json:"description"`
}

var IgnoredSettings = []string{
	"cluster.secret",
}

func ClusterSettingsFromRelease(release string) ([]ClusterSetting, error) {
	t, err := testserver.NewTestServer(
		testserver.CustomVersionOpt(release))
	if err != nil {
		return nil, err
	}
	pool, err := dbpgx.NewPoolFromUrl(t.PGURL().String())
	if err != nil {
		return nil, err
	}
	settings, err := GetLocalClusterSettings(pool)

	t.Stop()

	// Delete files
	files, err := filepath.Glob(filepath.Join(os.TempDir(), "cockroach-v*"))
	if err != nil {
		return settings, err
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return settings, err
		}
	}

	return settings, err
}

/*
// SaveClusterSettings save all cluster settings for all versions, but only if the combination of
// release, cpu and memory has not previously been run - otherwise it bails early
func SaveClusterSettings(url string) error {
	pool, err := dbpgx.NewPoolFromUrl(url)
	//if err != nil {
	//	return err
	//}
	//ds := NewDbDatasource(pool)

	rdb := releases.NewDbDatasource(pool)
	if err != nil {
		return err
	}
	rs, err := rdb.GetRecentReleaseNames(999999)
	if err != nil {
		return err
	}

	for _, r := range rs {
		err = SaveClusterSettingsForVersion(r, url)
		if err != nil {
			return err
		}
	}

	return nil
}

*/

/*
// SaveClusterSettingsForVersion saves all the cluster settings for a specific CRDB version, but only
// if the combination of release, cpu and memory has not been previously run - otherwise it bails early.
func SaveClusterSettingsForVersion(release string, url string) error {

	pool, err := dbpgx.NewPoolFromUrl(url)
	if err != nil {
		return err
	}
	ds := NewDbDatasource(pool)

	// Get host memory and CPU
	cpu := host.GetCpu()
	memoryBytes, err := host.GetMemory()
	if err != nil {
		return err
	}

	rs := make([]string, 0)
	if release == "all" {

	} else if strings.HasPrefix(release, "recent-") {
		rdb := releases.NewDbDatasource(pool)
		cntStr := strings.Replace(release, "recent-", "", 1)
		cnt, err := strconv.Atoi(cntStr)
		if err != nil {
			return err
		}
		rs, err = rdb.GetRecentReleaseNames(cnt)
		if err != nil {
			return err
		}
	} else {
		rs = append(rs, release)
	}

	// Iterate over releases
	for _, r := range rs {

		// Check to see if save run already exists, if it does, bail early - we've already captured the settings
		exists, err := ds.SaveRunExists(r, cpu, memoryBytes)
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

		err = ds.SaveRawSettings(rawSettings)
		if err != nil {
			return err
		}

		// Save the save run so we don't have to re-run later
		err = ds.SaveRun(r, cpu, memoryBytes)
		if err != nil {
			return err
		}
	}

	return nil

}

*/

// SummarizeSettings gets the raw settings and summarizes them into the settings_summary table
func SummarizeAndSaveSettings(url string) error {
	rsDs, err := NewDbDatasource(url)
	if err != nil {
		return err
	}
	rawSettings, err := rsDs.GetRawSettings()
	if err != nil {
		return err
	}

	rm, err := releases.NewReleasesManager(url)
	if err != nil {
		return err
	}

	rels, err := rm.GetReleases()
	if err != nil {
		return err
	}

	s := NewSummarizer(rawSettings, rels)

	summaries, err := s.Summarize()
	if err != nil {
		return err
	}

	return rsDs.SaveSettingsSummaries(summaries)
}
