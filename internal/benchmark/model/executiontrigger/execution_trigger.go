package executiontrigger

import (
	"context"
	"database/sql"
)

const (
	Command = "command"
	Cron    = "cron"
)

type ExecutionTrigger struct {
	ID   uint64
	Name string
}

func InitTable(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `
        INSERT INTO execution_trigger (name) 
        VALUES ($1), ($2) ON CONFLICT (name) DO NOTHING
    `)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.ExecContext(ctx, Command, Cron)
	if err != nil {
		return err
	}
	return nil
}

func Get(ctx context.Context, db *sql.DB, name string) (*ExecutionTrigger, error) {
	stmt, err := db.Prepare("SELECT id, name FROM execution_trigger WHERE name = $1")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	et := new(ExecutionTrigger)
	if err := stmt.QueryRowContext(ctx, name).Scan(&et.ID, &et.Name); err != nil {
		return nil, err
	}
	return et, nil
}
