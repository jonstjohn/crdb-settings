package gh

import (
	"context"
	"fmt"
	"github.com/google/go-github/v65/github"
	"time"
)

type Issue struct {
	ID        int64
	Number    int
	Title     string
	Url       string
	CreatedAt *time.Time
	ClosedAt  *time.Time
}

type Provider struct {
	Client *github.Client
}

func NewProvider(accessToken *string) Provider {
	if accessToken == nil || *accessToken == "" {
		client := github.NewClient(nil)
		return Provider{Client: client}

	}

	client := github.NewClient(nil).WithAuthToken(*accessToken)
	return Provider{Client: client}
}

func (p *Provider) SearchIssues(srch string) ([]Issue, error) {

	var issues []Issue
	q := fmt.Sprintf("%s repo:cockroachdb/cockroach is:pr", srch)
	//q := fmt.Sprintf("%s repo:jonstjohn/crdb-settings is:issue", srch)
	result, _, err := p.Client.Search.Issues(context.Background(), q, nil)
	if err != nil {
		return issues, err
	}

	// TODO - handle more than one page of results
	for _, i := range result.Issues {
		var closedAt *time.Time
		if i.ClosedAt != nil {
			closedAt = i.ClosedAt.GetTime()
		}
		issues = append(issues,
			Issue{Title: *i.Title, Url: *i.HTMLURL,
				ID: *i.ID, Number: *i.Number, CreatedAt: i.CreatedAt.GetTime(), ClosedAt: closedAt})
	}
	return issues, nil
}
