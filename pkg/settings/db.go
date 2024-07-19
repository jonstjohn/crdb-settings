package settings

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonstjohn/crdb-settings/pkg/dbpgx"
	"time"
)

type Db struct {
	Pool *pgxpool.Pool
}

/*
type SettingsSummaryRow struct {
	Variable           string
	Value              string
	Type               string
	Public             bool
	Description        string
	DefaultValue       string
	Origin             string
	Key                string
	FirstReleases      []string
	LastReleases       []string
	ValueChanges       []settings.Change
	DescriptionChanges []settings.Change
}

*/

const CreateRawTable = `
CREATE TABLE settings_raw (
    release_name string,
	cpu int,
	memory_bytes int,
	variable STRING NOT NULL,
	value STRING NOT NULL,
	type STRING NOT NULL,
	public BOOL NOT NULL,
	description STRING NOT NULL,
	default_value STRING NOT NULL,
	origin STRING NOT NULL,
	key STRING NOT NULL,
	updated TIMESTAMP NOT NULL DEFAULT now(),
	PRIMARY KEY (release_name, variable, cpu, memory_bytes),
	INDEX (variable, release_name)
)
`

const UpsertRaw = `
UPSERT INTO settings_raw (
	release_name, cpu, memory_bytes,
	variable, value, type,
	public, description, default_value,
	origin, key)
VALUES (
	$1, $2, $3,
	$4, $5, $6,
	$7, $8, $9,
	$10, $11
)
`

const CreateSaveRunTable = `
CREATE TABLE save_runs (
    release_name string,
	cpu int,
	memory_bytes int,
	updated TIMESTAMP NOT NULL DEFAULT now(),
	PRIMARY KEY (release_name, cpu, memory_bytes)
)
`

const UpsertSaveRun = `
UPSERT INTO save_runs (release_name, cpu, memory_bytes, updated)
VALUES ($1, $2, $3, now())
`

const OrderedRawSettingsSql = `
SELECT
	settings_raw.release_name,
	settings_raw.cpu,
	settings_raw.memory_bytes,
	settings_raw.variable,
	settings_raw.value,
	settings_raw.type,
	settings_raw.public,
	settings_raw.description,
	settings_raw.default_value,
	settings_raw.origin,
	settings_raw.key,
	settings_raw.updated
FROM
	settings_raw INNER JOIN releases ON settings_raw.release_name = releases.name
ORDER BY
	settings_raw.variable, releases.major, releases.minor, releases.patch,
	releases.beta_rc, releases.beta_rc_version, settings_raw.cpu,
	settings_raw.memory_bytes
`

const SelectSettingsForVersionSql = `
SELECT DISTINCT
	settings_raw.release_name,
	settings_raw.variable,
	settings_raw.value,
	settings_raw.type,
	settings_raw.public,
	settings_raw.description,
	settings_raw.default_value,
	settings_raw.origin,
	settings_raw.key,
	settings_raw.updated
FROM
	settings_raw
WHERE
	settings_raw.release_name = $1
ORDER BY public DESC, variable ASC
`

const CreateSummaryTable = `
CREATE TABLE settings_summary (
	variable STRING NOT NULL PRIMARY KEY,
	value STRING NOT NULL,
	type STRING NOT NULL,
	public BOOL NOT NULL,
	description STRING NOT NULL,
	default_value STRING NOT NULL,
	origin STRING NOT NULL,
	key STRING NOT NULL,
	first_releases STRING[] NULL,
	last_releases STRING[] NULL,
	value_changes JSONB NULL,
	description_changes JSONB NULL)
`

const UpsertSummarySql = `
UPSERT INTO settings_summary (
	variable, value, type, 
	public, description, default_value,
    origin, key, first_releases,
	last_releases, value_changes, description_changes)
VALUES (
	$1, $2, $3,
	$4, $5, $6,
	$7, $8, $9,
	$10, $11, $12)
`

const SelectSummarySql = `
SELECT variable, value, type, 
	public, description, default_value,
    origin, key, first_releases,
	last_releases, value_changes, description_changes
FROM settings_summary
ORDER BY variable
`

const CountSaveRun = `
SELECT count(*)
FROM save_runs
WHERE release_name = $1 AND cpu = $2 AND memory_bytes = $3
`

func NewDbDatasource(url string) (*Db, error) {
	pool, err := dbpgx.NewPoolFromUrl(url)
	if err != nil {
		return nil, err
	}
	return &Db{
		Pool: pool,
	}, nil
}

func (db *Db) GetSettingsForVersion(version string) (RawSettings, error) {
	rows, err := db.Pool.Query(context.Background(), SelectSettingsForVersionSql, version)
	if err != nil {
		return nil, err
	}

	sets := make([]RawSetting, 0)

	for rows.Next() {
		var releaseName string
		var variable string
		var value string
		var typ string
		var public bool
		var description string
		var defaultValue string
		var origin string
		var key string
		var updated time.Time
		err := rows.Scan(&releaseName, &variable, &value, &typ, &public,
			&description, &defaultValue, &origin, &key, &updated)
		if err != nil {
			return nil, err
		}
		sets = append(sets, RawSetting{
			ReleaseName:  releaseName,
			Variable:     variable,
			Value:        value,
			Type:         typ,
			Public:       public,
			Description:  description,
			DefaultValue: defaultValue,
			Origin:       origin,
			Key:          key,
			Updated:      updated,
		})
	}

	return sets, nil
}

func (db *Db) GetRawSettings() (RawSettings, error) {

	rows, err := db.Pool.Query(context.Background(), OrderedRawSettingsSql)
	if err != nil {
		return nil, err
	}

	sets := make([]RawSetting, 0)

	for rows.Next() {
		var releaseName string
		var cpu int
		var memoryBytes int64
		var variable string
		var value string
		var typ string
		var public bool
		var description string
		var defaultValue string
		var origin string
		var key string
		var updated time.Time
		err := rows.Scan(&releaseName, &cpu, &memoryBytes, &variable, &value, &typ, &public,
			&description, &defaultValue, &origin, &key, &updated)
		if err != nil {
			return nil, err
		}
		sets = append(sets, RawSetting{
			ReleaseName:  releaseName,
			Cpu:          cpu,
			MemoryBytes:  memoryBytes,
			Variable:     variable,
			Value:        value,
			Type:         typ,
			Public:       public,
			Description:  description,
			DefaultValue: defaultValue,
			Origin:       origin,
			Key:          key,
			Updated:      updated,
		})
	}

	return sets, nil
}

/*
func (db *Db) GetSettingSummaries() (settings.Summaries, error) {
	rows, err := db.Pool.Query(context.Background(), SelectSummarySql)
	if err != nil {
		return nil, err
	}

	summaries := make([]settings.Summary, 0)

	for rows.Next() {
		var variable string
		var value string
		var typ string
		var public bool
		var description string
		var defaultValue string
		var origin string
		var key string
		var firstReleases []string
		var lastReleases []string
		var valueChanges []settings.Change
		var descriptionChanges []settings.Change

		err := rows.Scan(&variable, &value, &typ, &public,
			&description, &defaultValue, &origin, &key,
			&firstReleases, &lastReleases,
			&value, &descriptionChanges)
		if err != nil {
			return nil, err
		}
		sets = append(sets, settings.RawSetting{
			ReleaseName:  releaseName,
			Cpu:          cpu,
			MemoryBytes:  memoryBytes,
			Variable:     variable,
			Value:        value,
			Type:         typ,
			Public:       public,
			Description:  description,
			DefaultValue: defaultValue,
			Origin:       origin,
			Key:          key,
		})
	}
	return nil
}

*/

func (db *Db) SaveRawSettings(rs RawSettings) error {
	for _, r := range rs {
		err := db.upsertRawSetting(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Db) SaveSettingsSummaries(ss Summaries) error {
	for _, s := range ss {
		err := db.upsertSummary(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Db) createRawTable() error {
	_, err := db.Pool.Exec(context.Background(), CreateRawTable)
	return err
}

func (db *Db) upsertRawSetting(r RawSetting) error {
	_, err := db.Pool.Exec(context.Background(), UpsertRaw,
		r.ReleaseName, r.Cpu, r.MemoryBytes,
		r.Variable, r.Value, r.Type,
		r.Public, r.Description, r.DefaultValue,
		r.Origin, r.Key,
	)
	return err
}

func (db *Db) createSaveRunTable() error {
	_, err := db.Pool.Exec(context.Background(), CreateSaveRunTable)
	return err
}

func (db *Db) SaveRunExists(releaseName string, cpu int, memoryBytes int64) (bool, error) {
	var cnt int
	err := db.Pool.QueryRow(context.Background(), CountSaveRun, releaseName, cpu, memoryBytes).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (db *Db) SaveRun(release string, cpu int, memory int64) error {
	_, err := db.Pool.Exec(context.Background(), UpsertSaveRun,
		release, cpu, memory)
	return err
}

func (db *Db) upsertSummary(summary Summary) error {

	valueChangesB, err := json.Marshal(summary.ValueChanges)
	descriptionChangesB, err := json.Marshal(summary.DescriptionChanges)
	if err != nil {
		return err
	}
	_, err = db.Pool.Exec(context.Background(), UpsertSummarySql,
		summary.Variable, summary.Value, summary.Type,
		summary.Public, summary.Description, summary.DefaultValue,
		summary.Origin, summary.Key, summary.FirstReleases,
		summary.LastReleases, valueChangesB, descriptionChangesB)
	return err
}

/*
UPSERT INTO settings_summary (
	variable, value, type,
	public, description, default_value,
    origin, key, first_releases, last_releases,
    value_changes, description_changes)
VALUES (
	$1, $2, $3,
	$4, $5, $6,
	$7, $8, $9,
	$10, $11)
*/
