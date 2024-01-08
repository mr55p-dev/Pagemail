package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var PageEventMap map[string]EventOutput[Page]

func init() {
	PageEventMap = make(map[string]EventOutput[Page])
}

type Client struct {
	// CreateRecord(*T) error
	// CreateRecords([]T) error
	// ReadRecordByField(field string, val any) (*T, error)
	// ReadRecordsByField(field string, val any) ([]T, error)
	// UpsertRecord(*T) error
	// UpsertRecords([]T) error
	// DeleteRecordsByField(field string, val any) error
	// AddListener() string
	// RemoveListener(id string)
	// DB() *sql.DB
	// Close()
	db *sqlx.DB
}

func (c *Client) Close() {
	c.db.Close()
}

func NewClient() *Client {
	conn := sqlx.MustOpen("sqlite3", "db/pagemail.sqlite3")

	return &Client{conn}
}
