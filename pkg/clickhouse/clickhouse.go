package clickhouse

import (
	"database/sql"
	"fmt"
	"hezzl/config"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type Clickhouse struct {
	db *sql.DB
}

func NewClickhouse(cfg config.ClickhouseConfig) (*sql.DB, error) {

	db, err := sql.Open("clickhouse", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error while trying to connect to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error while trying to ping db: %w", err)
	}

	return db, nil
}
