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
		client.log.WithError(err).ErrorCtx(c, "Failed creating user")
		sqlErr := err.(sqlite3.Error)
		switch sqlErr.Code {
		case sqlite3.ErrNo(sqlite3.ErrConstraintUnique), sqlite3.ErrNo(sqlite3.ErrConstraint):
			return errors.New("Invalid email")
		}
		return err
	}
	client.log.DebugCtx(c, "Created user")
	return nil
}

func (client *Client) ReadUserById(c context.Context, id string) (*User, error) {
	user := new(User)
	err := client.db.Get(user, `SELECT * FROM users WHERE id = ?`, id)
	if err != nil {
		client.log.ErrorCtx(c, "Failed to read user", err)
		return nil, err
	}
	client.log.DebugCtx(c, "Found user")
	return user, nil
}

func (client *Client) ReadUserByShortcutToken(c context.Context, token string) (*User, error) {
	user := new(User)
	err := client.db.GetContext(c, user, `SELECT * FROM users WHERE shortcut_token = ?`, token)
	if err != nil {
		client.log.ErrorCtx(c, "Failed to read user", "error", err)
		return nil, err
	}
	client.log.DebugCtx(c, "Found user")
	return user, nil
}

func (client *Client) ReadUserByEmail(c context.Context, email string) (*User, error) {
	user := new(User)
	err := client.db.GetContext(c, user, `SELECT * FROM users WHERE email = ?`, email)
	if err != nil {
		client.log.ErrorCtx(c, "Failed to read user", "email", email, "err", err.Error())
		return nil, err
	}
	client.log.DebugCtx(c, "Found user")
	return user, nil
}

func (client *Client) ReadUsersWithMail(c context.Context) (users []User, err error) {
	err = client.db.Select(&users, `SELECT * FROM users WHERE subscribed = true`)
	if err != nil {
		client.log.ErrorCtx(c, "Failed looking up user", err)
		return
	}
	client.log.DebugCtx(c, "Found users", "count", len(users))
	return
}

type UserTokenPair struct {
	UserId        string `db:"id"`
	ShortcutToken string `db:"shortcut_token"`
}

func (client *Client) ReadUserShortcutTokens(c context.Context) (out []UserTokenPair, err error) {
	err = client.db.Select(&out, `SELECT id, shortcut_token FROM users WHERE shortcut_token IS NOT NULL`)
	if err != nil {
		client.log.ErrorCtx(c, "Failed looking up users by token", err)
	}
	client.log.InfoCtx(c, "Found users with tokens")
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
		client.log.ErrorCtx(c, "Failed updating user", err)
		return err
	}
	client.log.InfoCtx(c, "Updated user")
	return nil
}
