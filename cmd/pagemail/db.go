package main

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
)

func getDb(ctx context.Context) *dbqueries.Queries {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	return dbqueries.New(conn)
}
