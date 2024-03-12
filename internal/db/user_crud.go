package db

import (
	"context"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3"
)

func (client *Client) CreateUser(c context.Context, u *User) error {
	_, err := client.db.NamedExec(`
		INSERT INTO users 
		VALUES (:id, :username, :email, :password, :name, :avatar, :subscribed, :shortcut_token, :has_readability, :created, :updated)`,
		u,
	)

	if err != nil {
		client.log.Errc(c, "Failed creating user", err)
		sqlErr := err.(sqlite3.Error)
		switch sqlErr.Code {
		case sqlite3.ErrNo(sqlite3.ErrConstraintUnique), sqlite3.ErrNo(sqlite3.ErrConstraint):
			return errors.New("Invalid email")
		}
		return err
	}
	client.log.DebugContext(c, "Created user")
	return nil
}

func (client *Client) ReadUserById(c context.Context, id string) (*User, error) {
	user := new(User)
	err := client.db.Get(user, `SELECT * FROM users WHERE id = ?`, id)
	if err != nil {
		client.log.Errc(c, "Failed to read user", err)
		return nil, err
	}
	client.log.DebugContext(c, "Found user")
	return user, nil
}

func (client *Client) ReadUserByShortcutToken(c context.Context, token string) (*User, error) {
	user := new(User)
	err := client.db.GetContext(c, user, `SELECT * FROM users WHERE shortcut_token = ?`, token)
	if err != nil {
		client.log.ErrorContext(c, "Failed to read user", "error", err)
		return nil, err
	}
	client.log.DebugContext(c, "Found user")
	return user, nil
}

func (client *Client) ReadUserByEmail(c context.Context, email string) (*User, error) {
	log := client.log.With("package", "db")
	user := new(User)
	err := client.db.GetContext(c, user, `SELECT * FROM users WHERE email = ?`, email)
	if err != nil {
		log.ErrorContext(c, "Failed to read user", "email", email, "err", err.Error())
		return nil, err
	}
	log.DebugContext(c, "Found user")
	return user, nil
}

func (client *Client) ReadUsersWithMail(c context.Context) (users []User, err error) {
	err = client.db.Select(&users, `SELECT * FROM users WHERE subscribed = true`)
	if err != nil {
		client.log.Errc(c, "Failed looking up user", err)
		return
	}
	client.log.DebugContext(c, "Found users", "count", len(users))
	return
}

type UserTokenPair struct {
	UserId        string `db:"id"`
	ShortcutToken string `db:"shortcut_token"`
}

func (client *Client) ReadUserShortcutTokens(c context.Context) (out []UserTokenPair, err error) {
	err = client.db.Select(&out, `SELECT id, shortcut_token FROM users WHERE shortcut_token IS NOT NULL`)
	if err != nil {
		client.log.Errc(c, "Failed looking up users by token", err)
	}
	client.log.InfoContext(c, "Found users with tokens")
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
		client.log.Errc(c, "Failed updating user", err)
		return err
	}
	client.log.InfoContext(c, "Updated user")
	return nil
}
