package dbpgx

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPoolFromUrl(url string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), url)
}
