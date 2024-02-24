package auth

import (
	"context"

	"github.com/mr55p-dev/pagemail/internal/db"
)

// maps tokens to user ids
var SESS_COOKIE string = "pm-auth-tkn"
var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type Authorizer interface {
	ValCredentialsAgainstUser(email, password string, user *db.User) bool
	ValUserAgainstPage(user *db.User, page *db.Page) bool
	GenPasswordHash(string) string
	GenSessionToken(user *db.User) string
	ValSessionToken(token string) (userId string)
	RevokeSessionToken(token string) (ok bool)
	GenShortcutToken(*db.User) string
}

type SecureAuthorizer struct {
	store map[string]string
}

func NewSecureAuthorizer(ctx context.Context, shortcutTokens ...db.UserTokenPair) Authorizer {
	baseStore := LoadShortcutTokens(shortcutTokens)
	return &SecureAuthorizer{
		store: baseStore,
	}
}

func LoadShortcutTokens(users []db.UserTokenPair) map[string]string {
	out := make(map[string]string, len(users))
	for _, v := range users {
		out[v.ShortcutToken] = v.UserId
	}
	return out
}

type TestAuthorizer struct { }

func NewTestAuthorizer() Authorizer {
	return &TestAuthorizer{}
}
