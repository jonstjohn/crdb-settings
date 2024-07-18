package releases

import (
	"regexp"
)

/*
Retrieves and parses remote version information via a Remote Provider
RemoteProvider: gets release information from Github
RemoteRelease: release parsed directly from remote yaml
*/
import (
	"bytes"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// releaseDataURL is the location of the YAML file maintained by the
// docs team where release information is encoded. This data is used
// to render the public CockroachDB releases page. We leverage the
// data in structured format to generate release information used
// for testing purposes.
const releaseDataURL = "https://raw.githubusercontent.com/cockroachdb/docs/main/src/current/_data/releases.yml"

// var namePattern = regexp.MustCompile(`^v(\d+).(\d+).(\d+)-?(beta|rc)?\.?(\d+)$`)
var namePattern = regexp.MustCompile(`^v(\d+).(\d+).(\d+)-?(beta|rc|alpha)?\.?(\d+)?$`)
var majorVersionPattern = regexp.MustCompile(`^v(\d+).(\d+)$`)

type Remote struct{}

func NewRemoteDataSource() *Remote {
	return &Remote{}
}

type CustomTime struct {
	time.Time
}

// RemoteRelease contains the information we extract from the YAML file in
// `releaseDataURL`.
type RemoteRelease struct {
	Name         string     `yaml:"release_name"`
	Withdrawn    bool       `yaml:"withdrawn"`
	CloudOnly    bool       `yaml:"cloud_only"`
	ReleaseType  string     `yaml:"release_type"` // Production or Testing
	ReleaseDate  CustomTime `yaml:"release_date"`
	MajorVersion string     `yaml:"major_version"` // e.g., v1.0, v23.1
}

type Version struct {
	Major         int
	Minor         int
	Patch         int
	BetaRc        string
	BetaRcVersion int
}

type ReleaseFilter struct {
	Widthdrawn  bool
	ReleaseType string
	From        time.Time
	To          time.Time
}

func (ct *CustomTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	layout := "2006-01-02"
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func ServerForVersion(v string) (testserver.TestServer, error) {
	return testserver.NewTestServer(
		testserver.CustomVersionOpt(v))
}

func (r *Remote) GetReleases() (Releases, error) {

	rels := Releases{}
	remoteReleases, err := r.GetRemoteReleases()
	if err != nil {
		return rels, nil
	}

	for _, rel := range remoteReleases {
		v := rel.Version()
		rels = append(rels, Release{
			Name:          rel.Name,
			Withdrawn:     rel.Withdrawn,
			CloudOnly:     rel.CloudOnly,
			ReleaseType:   rel.ReleaseType,
			ReleaseDate:   rel.ReleaseDate.Time,
			MajorVersion:  rel.MajorVersion,
			Major:         v.Major,
			Minor:         v.Minor,
			Patch:         v.Patch,
			BetaRc:        v.BetaRc,
			BetaRcVersion: v.BetaRcVersion,
		})
	}

	return rels, nil
}

func (r *Remote) GetRemoteReleases() ([]RemoteRelease, error) {
	resp, err := http.Get(releaseDataURL)
	if err != nil {
		return nil, fmt.Errorf("could not download release data: %w", err)
	}
	defer resp.Body.Close()

	var blob bytes.Buffer
	if _, err := io.Copy(&blob, resp.Body); err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var data []RemoteRelease
	if err := yaml.Unmarshal(blob.Bytes(), &data); err != nil { //nolint:yaml
		return nil, fmt.Errorf("failed to YAML parse release data: %w", err)
	}

	return data, nil
}

func (r *RemoteRelease) MajorMinorOnly() (int, int, error) {
	matches := majorVersionPattern.FindStringSubmatch((r.MajorVersion))
	major, err := strconv.Atoi(matches[1])
	minor, err := strconv.Atoi(matches[2])
	return major, minor, err

}

func (r *RemoteRelease) Version() Version {
	v := Version{}

	// Attempt to parse the name to get the major, minor, patch, etc
	matches := namePattern.FindStringSubmatch(r.Name)
	if matches != nil {
		major, err := strconv.Atoi(matches[1])
		minor, err := strconv.Atoi(matches[2])
		patch, err := strconv.Atoi(matches[3])
		if err != nil {
			return v
		}
		v.Major = major
		v.Minor = minor
		v.Patch = patch
		v.BetaRc = matches[4]
		if len(matches[5]) > 0 {
			betaRcVersion, err := strconv.Atoi(matches[5])
			if err != nil {
				return v
			}
			v.BetaRcVersion = betaRcVersion
		}
	} else { // if we can't parse the versions, use the major minor from major_version
		major, minor, _ := r.MajorMinorOnly()
		v.Major = major
		v.Minor = minor
		v.BetaRc = "beta" // TODO just call it beta for now
	}

	return v
}

func (r *Remote) GetReleaseSortedByVersion() ([]RemoteRelease, error) {
	releases, err := r.GetRemoteReleases()
	if err != nil {
		return nil, err
	}
	sort.Slice(releases, func(i int, j int) bool {
		iv := releases[i].Version()
		jv := releases[j].Version()
		if iv.Major < jv.Major {
			return true
		} else if iv.Major == jv.Major {
			if iv.Minor < jv.Minor {
				return true
			} else if iv.Minor == jv.Minor {
				if iv.Patch < jv.Patch {
					return true
				} else if iv.Patch == jv.Patch {
					if iv.BetaRc == "alpha" && (jv.BetaRc == "beta" || jv.BetaRc == "rc" || jv.BetaRc == "") {
						return true
					} else if iv.BetaRc == "beta" && (jv.BetaRc == "rc" || jv.BetaRc == "") {
						return true
					} else if iv.BetaRc == "rc" && jv.BetaRc == "" {
						return true
					} else if iv.BetaRc == jv.BetaRc {
						return iv.BetaRcVersion < jv.BetaRcVersion
					} else {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
		return true
	})
	return releases, nil
}

// Releases that we care about
// Latest major non-production release, if latest major is not in production
// All production releases that are still supported or within 6 months of no longer being supported
// Most recent major releases of everything else
/*
func GetReleasesWeCareAbout() ([]RemoteRelease, error) {
	releases, err := GetReleaseSortedByVersion()
	if err != nil {
		return nil, err
	}
	filtered := make([]RemoteRelease, 0)

	now := time.Now()
	for i, r := range releases {

		// Latest release and it isn't production, always include it
		if i == len(releases)-1 && r.ReleaseType != "production" {
			releases = append(releases, r)
			continue
		}

		// Is supported with some additional grace period
		if r.isSupportedWithGrace() {

		}
		if r.ReleaseDate.Time
	}
}

*/

func (r *Remote) GetReleasesByMajor() (map[string][]RemoteRelease, []string, error) {
	bymajors := make(map[string][]RemoteRelease)
	majors := make([]string, 0)

	releases, err := r.GetReleaseSortedByVersion()
	if err != nil {
		return bymajors, majors, err
	}

	for _, r := range releases {
		// Initialize releases array, if needed
		if ok := bymajors[r.MajorVersion]; ok == nil {
			majors = append(majors, r.MajorVersion)
			bymajors[r.MajorVersion] = make([]RemoteRelease, 0)
		}
		bymajors[r.MajorVersion] = append(bymajors[r.MajorVersion], r)
	}

	return bymajors, majors, nil
}

/*
func GetLatestReleases(numMajors int) ([]RemoteRelease, error) {
	releases, err := GetRemoteReleases()
	if err != nil {
		return nil, err
	}
	filtered := make([]RemoteRelease, 0)
	lastMajor := 0
	lastMinor := 0
	for _, r := range releases {
		major, minor := majorMinorFromReleaseName(r.Name)
		filtered = append(filtered, r)
	}
}

func majorMinorFromReleaseName(n string) (int, int) {

}
*/
