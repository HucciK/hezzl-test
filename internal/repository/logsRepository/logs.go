package logsRepository

import (
	"database/sql"
	"fmt"
	"hezzl/internal/core"
)

type LogsClickhouse struct {
	db *sql.DB
}

func NewLogsClickhouse(db *sql.DB) *LogsClickhouse {
	return &LogsClickhouse{
		db: db,
	}
}

func (l LogsClickhouse) ProcessBatch(batch []core.Item) error {
	tx, err := l.db.Begin()
	if err != nil {
		return fmt.Errorf("error while trying to begin transaction: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO items (Id, CampaignId, Name, Description, Priority, Removed, EventTime) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("error while trying to rollback transaction: %w", err)
		}
		return fmt.Errorf("error while trying to prepare sql stmt: %w", err)
	}

	for _, b := range batch {
		_, err := stmt.Exec(int32(b.Id), int32(b.CampaignId), b.Name, b.Description, int32(b.Priority), b.Removed, b.CreatedAt)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf("error while trying to rollback transaction: %w", err)
			}
			return fmt.Errorf("error while trying to exec sql stmt: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error while trying to commit transaction: %w", err)
	}

	return nil
}
