package auth

import "github.com/mr55p-dev/pagemail/pkg/db"

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
}

type SecureAuthorizer struct {
	store map[string]string
}

func NewSecureAuthorizer() Authorizer {
	return &SecureAuthorizer{
		store: make(map[string]string),
	}
}

type TestAuthorizer struct {
	userId string
}

func NewTestAuthorizer(userId string) Authorizer {
	return &TestAuthorizer{userId}
}
