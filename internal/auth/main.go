package auth

import (
	"context"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
)

// maps tokens to user ids
var SESS_COOKIE string = "pm-auth-tkn"
var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type Authorizer struct {
	store map[string]string
}

func NewAuthorizer(ctx context.Context) *Authorizer {
	return &Authorizer{
		store: make(map[string]string),
	}
}

func LoadShortcutTokens(users []dbqueries.User) map[string]string {
	out := make(map[string]string, len(users))
	for _, v := range users {
		out[v.ShortcutToken] = v.ID
	}
	return out
}
