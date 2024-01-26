package db

import (
	"log/slog"
	"time"

	"github.com/mr55p-dev/pagemail/pkg/tools"
)

type User struct {
	Id             string     `db:"id" log:"id" json:"id"`
	Username       string     `db:"username" log:"username" json:"username"`
	Email          string     `db:"email" log:"email" json:"email"`
	Password       string     `db:"password" json:"-"`
	Name           string     `db:"name" log:"name" json:"name"`
	Avatar         string     `db:"avatar" log:"avatar" json:"avatar"`
	Subscribed     bool       `db:"subscribed" log:"subscribed" json:"subscribed"`
	ShortcutToken  string     `db:"shortcut_token" log:"shortcut_token" json:"shortcut_token"`
	HasReadability bool       `db:"has_readability" log:"has_readability" json:"has_readability"`
	Created        *time.Time `db:"created" log:"created" json:"created"`
	Updated        *time.Time `db:"updated" log:"updated" json:"updated"`
}

func (user *User) LogValue() slog.Value {
	vals := tools.LogValue(user)
	return slog.GroupValue(vals...)
}

func NewUser(email, password string) *User {
	now := time.Now()
	uid := tools.GenerateNewId(10)
	token := tools.GenerateNewShortcutToken()
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
