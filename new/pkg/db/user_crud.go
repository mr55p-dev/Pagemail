package db

import (
	"context"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/mr55p-dev/pagemail/pkg/logging"
)

func (client *Client) CreateUser(c context.Context, u *User) error {
	_, err := client.db.NamedExec(`
		INSERT INTO users 
		VALUES (:id, :username, :email, :password, :name, :avatar, :subscribed, :shortcut_token, :has_readability, :created, :updated)`,
		u,
	)

	if err != nil {
		sqlErr := err.(sqlite3.Error)
		switch sqlErr.Code {
		case sqlite3.ErrNo(sqlite3.ErrConstraintUnique), sqlite3.ErrNo(sqlite3.ErrConstraint):
			return errors.New("Invalid email")
		}
		return err
	}

	return err
}

func (client *Client) ReadUserById(c context.Context, id string) (*User, error) {
	user := new(User)
	err := client.db.Get(user, `SELECT * FROM users WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return user, err
}
func (client *Client) ReadUserByEmail(c context.Context, email string) (*User, error) {
	user := new(User)
	err := client.db.Get(user, `SELECT * FROM users WHERE email = ?`, email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (client *Client) ReadUsersWithMail(c context.Context) (users []User, err error) {
	log.InfoContext(c, "Looking up all users with mail enabled")
	err = client.db.Select(users, `SELECT * FROM users`)
	if err != nil {
		log.ErrorContext(c, "Failed looking up user", logging.Error, err.Error())
		return
	}
	log.InfoContext(c, "Found users", "num-users", len(users))
	return
}
