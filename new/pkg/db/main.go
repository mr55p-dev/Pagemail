package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

var PageEventMap map[string]EventOutput[Page]

func init() {
	PageEventMap = make(map[string]EventOutput[Page])
}

type Client[T any] interface {
	TableName() string
	CreateRecord(*T) error
	// CreateRecords([]T) error
	// ReadRecordByField(field string, val any) (*T, error)
	// ReadRecordsByField(field string, val any) ([]T, error)
	// UpsertRecord(*T) error
	// UpsertRecords([]T) error
	// DeleteRecordsByField(field string, val any) error
	//
	// AddListener() string
	// RemoveListener(id string)
}

type Driver interface {
	DB() *sql.DB
	Close()
}

type DBDriver struct {
	db  *sql.DB
	log *slog.Logger
}

func (c *DBDriver) DB() *sql.DB {
	return c.db
}

func (c *DBDriver) Close() {
	c.db.Close()
}

func NewDriver(logger *slog.Logger) Driver {
	conn, err := sql.Open("sqlite3", "db/pagemail.sqlite3")
	if err != nil {
		panic(err)
	}

	return &DBDriver{conn, logger}
}
