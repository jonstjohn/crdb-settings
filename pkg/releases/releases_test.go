package releases

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServerForVersion(t *testing.T) {
	ts, err := ServerForVersion("v23.1.14")
	assert.Nil(t, err)
	pgurl := ts.PGURL()
	assert.NotNil(t, pgurl)

}

func TestGetReleases(t *testing.T) {
	rp := NewRemoteDataSource()
	rs, err := rp.GetRemoteReleases()
	assert.Nil(t, err)
	for _, r := range rs {
		assert.NotNil(t, r.ReleaseDate)
	}
}

func TestReleaseVersion(t *testing.T) {

	type NamePatternResult struct {
		Major         int
		Minor         int
		Patch         int
		BetaRc        string
		BetaRcVersion int
	}
	type NamePatternTest struct {
		Release RemoteRelease
		Result  NamePatternResult
	}
	tests := []NamePatternTest{
		NamePatternTest{
			Release: RemoteRelease{Name: "v23.1.14"},
			Result:  NamePatternResult{Major: 23, Minor: 1, Patch: 14},
		},
		NamePatternTest{
			Release: RemoteRelease{Name: "v23.2.0-beta.3"},
			Result:  NamePatternResult{Major: 23, Minor: 2, Patch: 0, BetaRc: "beta", BetaRcVersion: 3},
		},
		NamePatternTest{
			Release: RemoteRelease{Name: "v21.1.0-alpha.3"},
			Result:  NamePatternResult{Major: 21, Minor: 1, Patch: 0, BetaRc: "alpha", BetaRcVersion: 3},
		},
	}
	for _, n := range tests {
		v := n.Release.Version()
		assert.Equal(t, n.Result.Major, v.Major)
		assert.Equal(t, n.Result.Minor, v.Minor)
		assert.Equal(t, n.Result.Patch, v.Patch)

		if len(n.Result.BetaRc) > 0 {
			assert.Equal(t, n.Result.BetaRc, v.BetaRc)
		}

		if n.Result.BetaRcVersion != 0 {
			assert.Equal(t, n.Result.BetaRcVersion, v.BetaRcVersion)
		}
	}
}

func TestGetReleasesSortedByVersion(t *testing.T) {
	rp := NewRemoteDataSource()
	releases, err := rp.GetReleaseSortedByVersion()
	assert.Nil(t, err)
	m := make(map[int]map[int]map[int]RemoteRelease)
	for _, r := range releases {
		v := r.Version()
		if ok := m[v.Major]; ok == nil {
			m[v.Major] = make(map[int]map[int]RemoteRelease)
		}
		if ok := m[v.Major][v.Minor]; ok == nil {
			m[v.Major][v.Minor] = make(map[int]RemoteRelease)
		}
		m[v.Major][v.Minor][v.Patch] = r
	}
	assert.NotNil(t, nil)
}

func TestGetReleasesByMajor(t *testing.T) {
	rp := NewRemoteDataSource()
	releases, majors, err := rp.GetReleasesByMajor()
	assert.Nil(t, err)
	assert.NotNil(t, releases)
	assert.NotNil(t, majors)
}

func TestGet3MostRecentMajorReleases(t *testing.T) {
	rp := NewRemoteDataSource()
	releases, majors, err := rp.GetReleasesByMajor()
	assert.Nil(t, err)
	results := make([]RemoteRelease, 0)
	for _, m := range majors[len(majors)-3:] {
		// Get most recent not cloud only release
		for i := len(releases[m]) - 1; i >= 0; i-- {
			if !releases[m][i].CloudOnly && !releases[m][i].Withdrawn {
				results = append(results, releases[m][i])
				break
			}
		}
	}
	assert.NotEmpty(t, results)

}

func TestReleasesSortBy(t *testing.T) {
	rs := Releases{
		Release{
			Name:          "",
			Withdrawn:     false,
			CloudOnly:     false,
			ReleaseType:   "",
			ReleaseDate:   time.Now(),
			MajorVersion:  "",
			Major:         23,
			Minor:         1,
			Patch:         0,
			BetaRc:        "",
			BetaRcVersion: 0,
		},
		Release{
			Name:          "",
			Withdrawn:     false,
			CloudOnly:     false,
			ReleaseType:   "",
			ReleaseDate:   time.Now().Add(5 * time.Minute),
			MajorVersion:  "",
			Major:         23,
			Minor:         1,
			Patch:         0,
			BetaRc:        "rc",
			BetaRcVersion: 0,
		},
	}

	rs.SortBy(SortByVersion)
	assert.Equal(t, "rc", rs[0].BetaRc)

	rs.SortBy(SortByReleaseDate)
	assert.Equal(t, "", rs[0].BetaRc)
}
