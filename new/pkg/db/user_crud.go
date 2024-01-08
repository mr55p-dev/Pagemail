package db

import (
	"errors"

	"github.com/mattn/go-sqlite3"
)

func (c *DBDriver) CreateUser(u *User) error {
	c.log.Debug().Msgf("Inserting new user with id %s", u.Id)
	res, err := c.DB().Exec(`
		INSERT INTO users (
			id, username, email, password,
			name, avatar, subscribed, shortcut_token,
			has_readability, created, updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Id, u.Username, u.Email, u.Password, u.Name, u.Avatar,
		u.Subscribed, u.ShortcutToken, u.HasReadability, u.Created, u.Updated,
	)
	c.log.Debug().Msgf("Created new user with result %+v", res)

	if err != nil {
		c.log.Debug().Msgf("Created new user with err %s", err.Error())
		sqlErr := err.(sqlite3.Error)
		switch sqlErr.Code {
		case sqlite3.ErrNo(sqlite3.ErrConstraintUnique), sqlite3.ErrNo(sqlite3.ErrConstraint):
			c.log.Error().Msgf("Failed unique constraint check on %s", u.Email)
			return errors.New("Invalid email")
		}
	}
	return err
}

func (c *DBDriver) ReadUserById(id string) (*User, error) {
	c.log.Debug().Msgf("Checking user by id %s", id)
	row := c.DB().QueryRow(`SELECT * FROM users WHERE id = ?`, id)
	user := User{}
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Avatar,
		&user.Subscribed,
		&user.ShortcutToken,
		&user.HasReadability,
		&user.Created,
		&user.Updated,
	)
	c.log.Debug().Msgf("Loaded user %+v", user)
	return &user, err
}
func (c *DBDriver) ReadUserByEmail(email string) (*User, error) {
	c.log.Debug().Msgf("Checking user by email %s", email)
	row := c.DB().QueryRow(`SELECT * FROM users WHERE email = ?`, email)
	user := User{}
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Avatar,
		&user.Subscribed,
		&user.ShortcutToken,
		&user.HasReadability,
		&user.Created,
		&user.Updated,
	)
	c.log.Debug().Msgf("Loaded user %+v", user)
	return &user, err
}
