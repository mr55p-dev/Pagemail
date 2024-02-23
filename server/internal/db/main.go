package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mr55p-dev/pagemail/internal/logging"
)

type Client struct {
	log logging.Logger
	db  *sqlx.DB
}

func (c *Client) Close() {
	c.db.Close()
}

func NewClient(path string, log logging.Logger) *Client {
	conn := sqlx.MustOpen("sqlite3", path)
	return &Client{
		log: log,
		db:  conn,
	}
}
