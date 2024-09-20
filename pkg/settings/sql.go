package settings

/*
import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlExecutor struct {
	Pool *pgxpool.Pool
}

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
	settings_raw.updated,
FROM
	settings_raw INNER JOIN releases ON settings_raw.release_name = releases.name
ORDER BY
	settings_raw.variable, releases.major, releases.minor, releases.patch,
	releases.beta_rc, releases.beta_rc_version, settings_raw.cpu,
	settings_raw.memory_bytes
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
	first_releases []STRING NULL,
	last_releases []STRING NULL,
	value_changes JSONB NULL,
	description_changes JSONB NULL
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

const CountSaveRun = `
SELECT count(*)
FROM save_runs
WHERE release_name = $1 AND cpu = $2 AND memory_bytes = $3
`

func NewSqlExecutor(pool *pgxpool.Pool) *SqlExecutor {
	return &SqlExecutor{
		Pool: pool,
	}
}

func (s *SqlExecutor) CreateRawTable() error {
	_, err := s.Pool.Exec(context.Background(), CreateRawTable)
	return err
}

func (s *SqlExecutor) UpsertRawSetting(r RawSetting) error {
	_, err := s.Pool.Exec(context.Background(), UpsertRaw,
		r.ReleaseName, r.Cpu, r.MemoryBytes,
		r.Name, r.Value, r.Type,
		r.Public, r.Description, r.DefaultValue,
		r.Origin, r.Key,
	)
	return err
}

func (s *SqlExecutor) CreateSaveRunTable() error {
	_, err := s.Pool.Exec(context.Background(), CreateSaveRunTable)
	return err
}

func (s *SqlExecutor) SaveRunExists(releaseName string, cpu int, memoryBytes int64) (bool, error) {
	var cnt int
	err := s.Pool.QueryRow(context.Background(), CountSaveRun, releaseName, cpu, memoryBytes).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (s *SqlExecutor) UpsertSaveRun(release string, cpu int, memory int64) error {
	_, err := s.Pool.Exec(context.Background(), UpsertSaveRun,
		release, cpu, memory)
	return err
}

func (s *SqlExecutor) GetSettingsRawOrderedByVersion() ([]RawSetting, error) {

	rows, err := s.Pool.Query(context.Background(), OrderedRawSettingsSql)
	if err != nil {
		return nil, err
	}

	settings := make([]RawSetting, 0)

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
		err := rows.Scan(&releaseName, &cpu, &memoryBytes, &variable, &value, &typ, &public,
			&description, &defaultValue, &origin, &key)
		if err != nil {
			return nil, err
		}
		settings = append(settings, RawSetting{
			ReleaseName:  releaseName,
			Cpu:          cpu,
			MemoryBytes:  memoryBytes,
			Name:     variable,
			Value:        value,
			Type:         typ,
			Public:       public,
			Description:  description,
			DefaultValue: defaultValue,
			Origin:       origin,
			Key:          key,
		})
	}

	return settings, nil
}

func (s *SqlExecutor) UpsertSummary(summary Summary) error {

	valueChangesB, err := json.Marshal(summary.ValueChanges)
	descriptionChangesB, err := json.Marshal(summary.DescriptionChanges)
	if err != nil {
		return err
	}
	_, err = s.Pool.Exec(context.Background(), UpsertSummarySql,
		summary.Name, summary.Value, summary.Type,
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
