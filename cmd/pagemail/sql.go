package main

import (
	"context"
	"fmt"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// openDB establishes a connection to the given filepath and tests ping
func openDB(ctx context.Context, dsn string) (*pgx.Conn, error) {
	dbconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		// handle error
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to open connection: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping over connection: %w", err)
	}

	return conn, nil
}
