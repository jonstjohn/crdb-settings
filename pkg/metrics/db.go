package metrics

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

type RawRow struct {
	ReleaseName string
	Metric      string
	Help        string
	Type        string
	Updated     time.Time
}

const CreateRawTable = `
CREATE TABLE IF NOT EXISTS metrics_raw (
    release_name string,
	metric STRING NOT NULL,
	type STRING NOT NULL,
	help STRING NOT NULL,
	updated TIMESTAMP NOT NULL DEFAULT now(),
	PRIMARY KEY (release_name, metric),
	INDEX (variable, release_name)
)
`

type SaveRunsRow struct {
	ReleaseName string
	Updated     time.Time
}

const CreateSaveRunsTable = `
CREATE TABLE IF NOT EXISTS metrics_save_runs (
    release_name string,
	updated TIMESTAMP NOT NULL DEFAULT now(),
	PRIMARY KEY (release_name)
)
`

const UpsertRaw = `
UPSERT INTO metrics_raw (release_name, metric, type, help) VALUES ($1, $2, $3, $4)
`

const UpsertSaveRun = `
UPSERT INTO metrics_save_runs (release_name) VALUES ($1, $2)
`

const SelectMetricsForReleaseSql = `
SELECT release_name, metric, type, help, updated
FROM metrics_raw 
WHERE release_name = $1
ORDER BY metric ASC
`

const SelectSaveRunsForReleaseSql = `
SELECT release_name, updated
FROM metrics_save_runs
WHERE release_name = $1
`

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

func (db *Db) createRawTableIfNotExists() error {
	_, err := db.Pool.Exec(context.Background(), CreateRawTable)
	return err
}

func (db *Db) createSaveRunsTableIfNotExists() error {
	_, err := db.Pool.Exec(context.Background(), CreateSaveRunsTable)
	return err
}

func (db *Db) UpsertRaw(releaseName, metric Metric) error {
	_, err := db.Pool.Exec(context.Background(), UpsertRaw,
		releaseName, metric.Name, metric.Type, metric.Help,
	)
	return err
}

func (db *Db) UpsertSaveRun(releaseName, metric Metric) error {
	_, err := db.Pool.Exec(context.Background(), UpsertSaveRun, releaseName)
	return err
}

func (db *Db) SelectRaw(releaseName string) ([]RawRow, error) {
	rows, err := db.Pool.Query(context.Background(), SelectMetricsForReleaseSql, releaseName)
	if err != nil {
		return nil, err
	}

	rs := make([]RawRow, 0)

	for rows.Next() {

		var releaseName string
		var metric string
		var typ string
		var help string
		var updated time.Time
		err := rows.Scan(&releaseName, &metric, &typ, &help, &updated)
		if err != nil {
			return nil, err
		}
		rs = append(rs, RawRow{
			ReleaseName: releaseName, Metric: metric,
			Type: typ, Help: help, Updated: updated,
		})
	}

	return rs, nil
}

func (db *Db) SelectSaveRuns(releaseName string) ([]SaveRunsRow, error) {

	rows, err := db.Pool.Query(context.Background(), SelectSaveRunsForReleaseSql, releaseName)
	if err != nil {
		return nil, err
	}

	rs := make([]SaveRunsRow, 0)

	for rows.Next() {

		var releaseName string
		var updated time.Time
		err := rows.Scan(&releaseName, &updated)
		if err != nil {
			return nil, err
		}
		rs = append(rs, SaveRunsRow{
			ReleaseName: releaseName, Updated: updated,
		})
	}

	return rs, nil
}
