package dbqueries

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mr55p-dev/pagemail/db"
)

func MustGetDB(ctx context.Context, path string) *sql.DB {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	err = conn.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	return conn
}

func LoadTables(ctx context.Context, conn *sql.DB) error {
	_, err := conn.ExecContext(ctx, db.Schema)
	if err != nil {
		return err
	}
	return nil
}
