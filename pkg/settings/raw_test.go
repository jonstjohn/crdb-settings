package settings

import (
	"blatta/pkg/releases"
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func releasesFromFile() ([]releases.Release, error) {
	file, err := os.Open("testdata/releases")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	rs := make([]releases.Release, 0)
	for scanner.Scan() {
		l := scanner.Text()
		parts := strings.Split(l, "||")
		name := parts[0]
		withdrawn := parts[1] == "true"
		releaseType := parts[2]
		releaseDate, err := time.Parse("2006-01-02 15:04:05", parts[3])
		if err != nil {
			return nil, err
		}
		majorVersion := parts[4]
		major, err := strconv.Atoi(parts[5])
		if err != nil {
			return nil, err
		}
		minor, err := strconv.Atoi(parts[6])
		if err != nil {
			return nil, err
		}
		patch, err := strconv.Atoi(parts[7])
		if err != nil {
			return nil, err
		}
		betaRc := parts[8]
		betaRcVersion, err := strconv.Atoi(parts[9])
		if err != nil {
			return nil, err
		}
		rs = append(rs, releases.Release{
			Name:          name,
			Withdrawn:     withdrawn,
			ReleaseType:   releaseType,
			ReleaseDate:   releaseDate,
			MajorVersion:  majorVersion,
			Major:         major,
			Minor:         minor,
			Patch:         patch,
			BetaRc:        betaRc,
			BetaRcVersion: betaRcVersion,
		})

	}

	return rs, nil
}

func rawSettingsFromFile(names ...string) (*RawSettings, error) {

	rs := make([]RawSetting, 0)
	for _, name := range names {
		file, err := os.Open(fmt.Sprintf("testdata/%s", name))
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			l := scanner.Text()
			parts := strings.Split(l, "||")
			cpu, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}
			rs = append(rs, RawSetting{ReleaseName: parts[0], Cpu: cpu,
				Variable: parts[2], Value: parts[3], Description: parts[4]})

		}
	}

	return (*RawSettings)(&rs), nil
}

func TestRawSettingsMetaForVariable(t *testing.T) {

	rels, err := releasesFromFile()
	assert.Nil(t, err)

	recoveryVariable := "kv.snapshot_recovery.max_rate"
	runnersVariable := "sql.distsql.num_runners"
	rawSettings, err := rawSettingsFromFile(recoveryVariable, runnersVariable)
	assert.Nil(t, err)

	recoveryMeta := rawSettings.MetaForVariable(recoveryVariable, rels)
	assert.Equal(t, false, recoveryMeta.hostDependent)

	assert.Equal(t, 4, len(recoveryMeta.firstReleases))
	assert.Equal(t, "v23.1.15", recoveryMeta.mostRecent.ReleaseName)
	assert.Equal(t, 1, len(recoveryMeta.valueChanges))
	assert.Equal(t, "v22.1.10", recoveryMeta.valueChanges[0].Release)

	runnersMeta := rawSettings.MetaForVariable(runnersVariable, rels)
	assert.Equal(t, true, runnersMeta.hostDependent)

}
