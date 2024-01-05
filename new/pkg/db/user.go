package db

import (
	"time"

	"github.com/mr55p-dev/pagemail/pkg/tools"
)

type User struct {
	Id             string     `db:"id"`
	Username       string     `db:"username"`
	Email          string     `db:"email"`
	Password       string     `db:"password"`
	Name           string     `db:"name"`
	Avatar         string     `db:"avatar"`
	Subscribed     bool       `db:"subscribed"`
	ShortcutToken  string     `db:"shortcutoken"`
	HasReadability bool       `db:"has_readability"`
	Created        *time.Time `db:"created"`
	Updated        *time.Time `db:"updated"`
}

func NewUser(email, password string) *User {
	now := time.Now()
	uid := tools.GenerateNewId(20)
	token := tools.GenerateNewShortcutToken(uid)
	return &User{
		Id:             uid,
		Email:          email,
		Password:       password,
		ShortcutToken:  token,
		HasReadability: false,
		Created:        &now,
		Updated:        &now,
	}
}
