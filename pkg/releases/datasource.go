package releases

type Provider interface {
	GetReleases() (Releases, error)
}

type Persister interface {
	SaveReleases(Releases) error
}
