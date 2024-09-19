package gh

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
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

func (m *Manager) GetIssuesForSetting(setting string) ([]Issue, error) {
	return m.Provider.SearchIssues(setting)
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

	for _, s := range settings {

		logrus.Info(fmt.Sprintf("Processing setting '%s'", s))

		issues, err := m.GetIssuesForSetting(s)
		if err != nil {
			return err
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
