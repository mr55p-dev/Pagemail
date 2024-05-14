package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(ctx context.Context, path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to driver: %w", err)
	}
	err = conn.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error pinging over connection: %w", err)
	}
	return conn, nil
}

func LoadSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return err
	}
	return nil
}
