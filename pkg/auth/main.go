package auth

import (
	"context"

	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/logging"
)

// maps tokens to user ids
var SESS_COOKIE string = "pm-auth-tkn"
var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
var log logging.Log

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

func init() {
	log = logging.GetLogger("authorizer")
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

type TestAuthorizer struct {
	userId string
}

func NewTestAuthorizer(userId string) Authorizer {
	return &TestAuthorizer{userId}
}
