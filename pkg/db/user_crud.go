package db

import (
	"context"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/mr55p-dev/pagemail/pkg/logging"
)

func (client *Client) CreateUser(c context.Context, u *User) error {
	res, err := client.db.NamedExec(`
		INSERT INTO users 
		VALUES (:id, :username, :email, :password, :name, :avatar, :subscribed, :shortcut_token, :has_readability, :created, :updated)`,
		u,
	)

	if err != nil {
		log.ErrContext(c, "Failed creating user", err, logging.UserId, u.Id)
		sqlErr := err.(sqlite3.Error)
		switch sqlErr.Code {
		case sqlite3.ErrNo(sqlite3.ErrConstraintUnique), sqlite3.ErrNo(sqlite3.ErrConstraint):
			return errors.New("Invalid email")
		}
		return err
	}
	rows, _ := res.RowsAffected()
	log.DebugContext(c, "Created user", logging.User, u, logging.Rows, rows)
	return nil
}

func (client *Client) ReadUserById(c context.Context, id string) (*User, error) {
	user := new(User)
	err := client.db.Get(user, `SELECT * FROM users WHERE id = ?`, id)
	if err != nil {
		log.ErrContext(c, "Failed to read user", err, logging.UserId, id)
		return nil, err
	}
	log.DebugContext(c, "Found user", logging.User, user)
	return user, nil
}
func (client *Client) ReadUserByEmail(c context.Context, email string) (*User, error) {
	user := new(User)
	err := client.db.GetContext(c, user, `SELECT * FROM users WHERE email = ?`, email)
	if err != nil {
		log.ErrContext(c, "Failed to read user", err, logging.UserMail, email)
		return nil, err
	}
	log.DebugContext(c, "Found user", logging.User, user)
	return user, nil
}

func (client *Client) ReadUsersWithMail(c context.Context) (users []User, err error) {
	err = client.db.Select(&users, `SELECT * FROM users`)
	if err != nil {
		log.ErrContext(c, "Failed looking up user", err)
		return
	}
	log.DebugContext(c, "Found users", "count", len(users))
	return
}
