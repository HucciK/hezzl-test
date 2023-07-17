package postgres

import (
	"database/sql"
	"fmt"
	"hezzl/config"

	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error while trying to open db connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error while trying to ping db: %w", err)
	}

	return db, nil
}
