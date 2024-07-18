package settings

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

// Cluster settings (e.g., from crdb_internal)

type ClusterSetting struct {
	Variable     string
	Value        string
	Type         string
	Public       bool
	Description  string
	DefaultValue string
	Origin       string
	Key          string
}

const ColCountSql = "SELECT count(*) FROM information_schema.columns WHERE table_name = 'cluster_settings' AND table_schema = 'crdb_internal'"

// GetLocalClusterSettings gets settings from the local cluster. It is necessary to run it on one of the cluster nodes,
// either using a regular cluster or a testing cluster, since we also capture CPU and memory information from the
// cluster. Because columns have been added to crdb_internal.cluster_settings over time, this dynamically determines
// what columns to cpature data from and uses nil or default value for the rest
func GetLocalClusterSettings(pool *pgxpool.Pool) ([]ClusterSetting, error) {
	var cnt int
	pool.QueryRow(context.Background(), ColCountSql).Scan(&cnt)

	// Currently, only 8 columns are supported. If more columns are added then these will need to be handled
	if cnt > 8 {
		cnt = 8
	}

	columns := []string{"variable", "value", "type", "public", "description", "default_value", "origin", "key"}

	sql := fmt.Sprintf("SELECT %s FROM crdb_internal.cluster_settings", strings.Join(columns[0:cnt], ","))
	rows, err := pool.Query(context.Background(), sql) // , targets[0:cnt]...)
	if err != nil {
		return nil, err
	}

	settings := make([]ClusterSetting, 0)
	for rows.Next() {
		var variable string
		var value string
		var typ string
		var public bool
		var description string
		var defaultValue string
		var origin string
		var key string
		targets := []interface{}{
			&variable, &value, &typ, &public, &description, &defaultValue, &origin, &key,
		}
		err := rows.Scan(targets[0:cnt]...)
		if err != nil {
			return settings, err
		}
		settings = append(settings, ClusterSetting{
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
	return settings, nil

}
