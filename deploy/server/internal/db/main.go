package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var PageEventMap map[string]EventOutput[Page]
var log logging.Log

func init() {
	log = logging.GetLogger("db")
	PageEventMap = make(map[string]EventOutput[Page])
}

type Client struct {
	db *sqlx.DB
}

func (c *Client) Close() {
	c.db.Close()
}

func NewClient(path string) *Client {
	conn := sqlx.MustOpen("sqlite3", path)

	return &Client{conn}
}
