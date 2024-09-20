package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// openDB establishes a connection to the given filepath and tests ping
func openDB(ctx context.Context, path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", "file:"+path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open connection: %w", err)
	}

	err = conn.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping over connection: %w", err)
	}

	return conn, nil
}
