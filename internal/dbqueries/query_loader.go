package dbqueries

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func MustGetQueries(ctx context.Context, path string) (*Queries, func() error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	err = conn.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	return New(conn), conn.Close
}
