package gh

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v65/github"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	Provider Provider
	Db       *Db
}

func NewManager(accessToken *string, url string) (*Manager, error) {
	db, err := NewDbDatasource(url)
	if err != nil {
		return nil, err
	}
	provider := NewProvider(accessToken)
	return &Manager{Provider: provider, Db: db}, err
}

func (m *Manager) SearchIssuesForSetting(setting string) ([]Issue, error) {
	return m.Provider.SearchIssues(setting)
}

func (m *Manager) GetIssuesForSetting(setting string) ([]Issue, error) {
	rows, err := m.Db.GetIssuesForSetting(setting)
	if err != nil {
		return nil, err
	}

	issues := make([]Issue, len(rows))
	for i, row := range rows {
		issues[i] = Issue{
			ID:        row.Id,
			Number:    row.Number,
			Title:     row.Title,
			Url:       row.Url,
			CreatedAt: row.Created,
			ClosedAt:  row.Closed,
		}

	}
	return issues, nil
}

func (m *Manager) UpdateIssuesForSetting(setting string) error {

	var settings []string
	var err error
	if setting == "all" || strings.HasPrefix(setting, "oldest-") {
		cnt := math.MaxInt
		if strings.HasPrefix(setting, "oldest-") {
			cntStr := strings.Replace(setting, "oldest-", "", 1)
			cnt, err = strconv.Atoi(cntStr)
			if err != nil {
				return err
			}
		}
		s, err := m.GetOldestSettingStrings(cnt)
		if err != nil {
			return err
		}
		settings = append(settings, s...)
	} else {
		settings = append(settings, setting)
	}

	// Rate limiter: 10 requests per minute = 1 request every 6 seconds
	rateLimiter := time.NewTicker(6 * time.Second)
	defer rateLimiter.Stop()

	for i, s := range settings {

		// Wait for rate limiter, except for the first request
		if i > 0 {
			<-rateLimiter.C
		}

		logrus.Info(fmt.Sprintf("Processing setting '%s'", s))

		issues, err := m.SearchIssuesForSetting(s)
		if err != nil {
			// Check if it's a rate limit error
			if rateLimitErr, ok := err.(*github.RateLimitError); ok {
				// Calculate wait time until rate limit resets
				waitDuration := time.Until(rateLimitErr.Rate.Reset.Time)
				// Add a small buffer (1 second) to ensure the limit has reset
				waitDuration += time.Second

				logrus.Warnf("Rate limit exceeded, waiting %v until reset at %v",
					waitDuration.Round(time.Second), rateLimitErr.Rate.Reset.Time)
				time.Sleep(waitDuration)

				// Retry the request
				issues, err = m.SearchIssuesForSetting(s)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		for _, i := range issues {
			err := m.Db.SaveSettingIssue(s, i)
			if err != nil {
				return err
			}
		}

		err = m.Db.UpdateSettingProcessed(s)
		if err != nil {
			return err
		}

		logrus.Info(fmt.Sprintf("Updated setting '%s' with %d issues", s, len(issues)))
	}

	return nil
}

func (m *Manager) GetOldestSettingStrings(cnt int) ([]string, error) {
	return m.Db.GetOldestSettingStrings(cnt)
}
