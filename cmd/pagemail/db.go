package main

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
)

func mustGetDb(ctx context.Context, path string) (*dbqueries.Queries, func() error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	err = conn.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	return dbqueries.New(conn), conn.Close
}
