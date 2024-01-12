package db

import (
	"context"
	"errors"
	"time"

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

type UserTokenPair struct {
	UserId        string `db:"id"`
	ShortcutToken string `db:"shortcut_token"`
}

func (client *Client) ReadUserShortcutTokens(c context.Context) (out []UserTokenPair, err error) {
	err = client.db.Select(&out, `SELECT id, shortcut_token FROM users WHERE shortcut_token IS NOT NULL`)
	if err != nil {
		log.ErrContext(c, "Failed looking up users by token", err)
	}
	log.InfoContext(c, "Found users with tokens", logging.Rows, len(out))
	return
}

func (client *Client) UpdateUser(c context.Context, user *User) error {
	now := time.Now()
	user.Updated = &now
	_, err := client.db.NamedExecContext(c, `
		UPDATE users SET 
			username = :username,
			email = :email,
			password = :password,
			name = :name,
			avatar = :avatar,
			subscribed = :subscribed,
			shortcut_token = :shortcut_token,
			has_readability = :has_readability,
			updated = :updated
		WHERE id = :id
	`, user)
	if err != nil {
		log.ErrContext(c, "Failed updating user", err, logging.User, user)
		return err
	}
	log.InfoContext(c, "Updated user", logging.User, user)
	return nil
}
