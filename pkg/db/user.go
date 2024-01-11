package db

import (
	"log/slog"
	"time"

	"github.com/mr55p-dev/pagemail/pkg/tools"
)

type User struct {
	Id             string     `db:"id" log:"id"`
	Username       string     `db:"username" log:"username"`
	Email          string     `db:"email" log:"email"`
	Password       string     `db:"password"`
	Name           string     `db:"name" log:"name"`
	Avatar         string     `db:"avatar" log:"avatar"`
	Subscribed     bool       `db:"subscribed" log:"subscribed"`
	ShortcutToken  string     `db:"shortcut_token" log:"shortcut_token"`
	HasReadability bool       `db:"has_readability" log:"has_readability"`
	Created        *time.Time `db:"created" log:"created"`
	Updated        *time.Time `db:"updated" log:"updated"`
}

func (user *User) LogValue() slog.Value {
	vals := tools.LogValue(user)
	return slog.GroupValue(vals...)
}

func NewUser(email, password string) *User {
	now := time.Now()
	uid := tools.GenerateNewId(10)
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
