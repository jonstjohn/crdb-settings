package releases

type Manager struct {
	Db *Db
}

func NewReleasesManager(url string) (*Manager, error) {

	db, err := NewDbDatasource(url)
	if err != nil {
		return nil, err
	}
	return &Manager{Db: db}, err
}

func (rm *Manager) GetReleases() (Releases, error) {
	rows, err := rm.Db.getReleasesRows()
	if err != nil {
		return nil, err
	}

	rels := make([]Release, len(rows))
	for i, r := range rows {
		rels[i] = Release{
			Name:          r.Name,
			Withdrawn:     r.Withdrawn,
			CloudOnly:     r.CloudOnly,
			ReleaseType:   r.ReleaseType,
			ReleaseDate:   r.ReleaseDate,
			MajorVersion:  r.MajorVersion,
			Major:         r.Major,
			Minor:         r.Minor,
			Patch:         r.Patch,
			BetaRc:        r.BetaRc,
			BetaRcVersion: r.BetaRcVersion,
		}
	}
	return rels, nil
}

func (rm *Manager) UpdateReleases() error {
	remote := NewRemoteDataSource()
	releasesFromRemote, err := remote.GetReleases()
	if err != nil {
		return err
	}
	return rm.Db.SaveReleases(releasesFromRemote)
}

func (rm *Manager) GetRecentReleaseNames(cnt int) ([]string, error) {
	return rm.Db.GetRecentReleaseNames(cnt)
}
