package db

import (
	"context"
	"errors"

	"github.com/mattn/go-sqlite3"
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
