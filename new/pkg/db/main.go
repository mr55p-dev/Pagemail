package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
)

var PageEventMap map[string]EventOutput[Page]

func init() {
	PageEventMap = make(map[string]EventOutput[Page])
}

type AbsClient interface {
	DB() *sql.DB
	Close()
	CreateUser(*User) error
	ReadUserById(string) (*User, error)
	ReadUserByEmail(string) (*User, error)

	CreatePage(*Page) error
	UpsertPage(*Page) error

	ReadPagesByUserId(string) ([]Page, error)
	DeletePagesByUserId(string) (int, error)

	AddPageListener(id string, output EventOutput[Page])
	RemovePageListener(id string)
}

type Client struct {
	db  *sql.DB
	log zerolog.Logger
}

func (c *Client) DB() *sql.DB {
	return c.db
}

func (c *Client) Close() {
	c.db.Close()
}

func NewClient(logger zerolog.Logger) AbsClient {
	conn, err := sql.Open("sqlite3", "db/pagemail.sqlite3")
	if err != nil {
		panic(err)
	}

	return &Client{conn, logger}
}
