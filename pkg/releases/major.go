package releases

type MajorVersion struct {
	MajorVersion           string
	Releases               Releases
	FirstTestingRelease    *Release
	FirstProductionRelease *Release
	LastTestingRelease     *Release
	LastProductionRelease  *Release
}

type MajorVersionSummary struct {
	MajorVersions []MajorVersion
	LatestRelease *Release
}

func NewMajorVersionSummaryFromReleases(rs *Releases) *MajorVersionSummary {
	rs.SortBy(SortByVersion)

	all := make(map[string][]Release)
	testing := make(map[string][]Release)
	production := make(map[string][]Release)
	m := make(map[string]bool)
	mvs := make([]string, 0)
	for _, r := range *rs {

		// Store every major version
		if !m[r.MajorVersion] {
			m[r.MajorVersion] = true
			mvs = append(mvs, r.MajorVersion)
		}

		all[r.MajorVersion] = append(all[r.MajorVersion], r)

		// Testing releases
		if r.ReleaseType == "Testing" {
			if _, ok := testing[r.MajorVersion]; !ok {
				testing[r.MajorVersion] = make([]Release, 0)
			}
			testing[r.MajorVersion] = append(testing[r.MajorVersion], r)
		}

		// Production release
		if r.ReleaseType == "Production" {
			if _, ok := production[r.MajorVersion]; !ok {
				production[r.MajorVersion] = make([]Release, 0)
			}
			production[r.MajorVersion] = append(production[r.MajorVersion], r)
		}
	}

	summary := MajorVersionSummary{MajorVersions: make([]MajorVersion, len(mvs))}

	for i, mv := range mvs {
		summary.MajorVersions[i] = MajorVersion{
			MajorVersion:           mv,
			Releases:               all[mv],
			FirstTestingRelease:    &testing[mv][0],
			FirstProductionRelease: &production[mv][0],
			LastTestingRelease:     &testing[mv][len(testing[mv])-1],
			LastProductionRelease:  &production[mv][len(production[mv])-1],
		}
	}

	return &summary
}

/*
func GetMajorVersionSummary(releases Releases) *MajorVersionSummary {
	releases.SortBy(SortByVersion)
	mvs := MajorVersionSummary{MajorVersions: make([]MajorVersion, 0)}
	for _, m := range releases.MajorVersions() {
		mvs.MajorVersions = append(mvs.MajorVersions, MajorVersion{
			MajorVersion:           m,
			FirstTestingRelease:    releases.FirstTestingReleaseForMajorVersion(m),
			FirstProductionRelease: releases.FirstProductionReleaseForMajorVersion(m),
			LastTestingRelease: releases
		})
		mvs.LatestRelease = releases.LatestReleaseForMajorVersion(m)

	}
	return &mvs
}

*/
