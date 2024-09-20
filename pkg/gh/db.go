package gh

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonstjohn/crdb-settings/pkg/dbpgx"
	"time"
)

type Db struct {
	Url  string
	Pool *pgxpool.Pool
}

type SettingsGithubIssuesRow struct {
	Variable  string
	Id        int64
	Number    int
	Title     string
	Url       string
	Processed *time.Time
	Closed    *time.Time
	Created   *time.Time
}

func NewDbDatasource(url string) (*Db, error) {
	pool, err := dbpgx.NewPoolFromUrl(url)
	if err != nil {
		return nil, err
	}
	return &Db{
		Url:  url,
		Pool: pool,
	}, nil
}

func (db *Db) SaveSettingIssue(setting string, issue Issue) error {
	sql := "UPSERT INTO settings_github_issues (variable, id, number, url, title, created, closed, processed) VALUES ($1, $2, $3, $4, $5, $6, $7, now())"
	_, err := db.Pool.Exec(context.Background(), sql, setting, issue.ID, issue.Number, issue.Url, issue.Title, issue.CreatedAt, issue.ClosedAt)
	return err
}

func (db *Db) UpdateSettingProcessed(setting string) error {
	sql := "UPSERT INTO settings_github_processed (variable, processed) values ($1, now())"
	_, err := db.Pool.Exec(context.Background(), sql, setting)
	return err
}

func (db *Db) GetOldestSettingStrings(cnt int) ([]string, error) {
	sql := "WITH sr AS (SELECT variable FROM settings_raw GROUP BY variable) SELECT sr.variable FROM sr LEFT JOIN settings_github_processed sgp ON sr.variable = sgp.variable ORDER BY sgp.processed, sr.variable ASC LIMIT $1"
	rows, err := db.Pool.Query(context.Background(), sql, cnt)
	if err != nil {
		return nil, err
	}

	var strs []string

	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		strs = append(strs, s)
	}
	return strs, nil

}

func (db *Db) GetIssuesForSetting(setting string) ([]SettingsGithubIssuesRow, error) {
	sql := "SELECT variable, id, number, title, url, processed, closed, created FROM settings_github_issues WHERE variable = $1" +
		"ORDER BY closed DESC"
	rows, err := db.Pool.Query(context.Background(), sql, setting)
	if err != nil {
		return nil, err
	}
	var issues []SettingsGithubIssuesRow

	for rows.Next() {
		var issue SettingsGithubIssuesRow
		err := rows.Scan(&issue.Variable, &issue.Id, &issue.Number, &issue.Title, &issue.Url, &issue.Processed, &issue.Closed, &issue.Created)
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}

	return issues, nil
}
