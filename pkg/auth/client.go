package auth

import (
	"crypto/sha256"

	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/tools"
)

// Validation logic
func (*SecureAuthorizer) ValUserAgainstPage(user *db.User, page *db.Page) bool {
	return user.Id == page.UserId
}

func (a *SecureAuthorizer) ValCredentialsAgainstUser(email, password string, user *db.User) (isValid bool) {
	emailValid := email == user.Email
	passwordValid := user.Password == a.GenPasswordHash(password)

	return emailValid && passwordValid
}

func (*SecureAuthorizer) GenPasswordHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return string(h.Sum(nil))
}

// Session tokens
func (a *SecureAuthorizer) GenSessionToken(user *db.User) string {
	tkn := tools.GenerateNewId(50)
	a.store[tkn] = user.Id
	return tkn
}

func (a *SecureAuthorizer) ValSessionToken(token string) string {
	uid := a.store[token]
	return uid
}

func (a *SecureAuthorizer) RevokeSessionToken(token string) bool {
	_, ok := a.store[token]
	delete(a.store, token)
	return ok
}

func (a *SecureAuthorizer) GenShortcutToken(user *db.User) string {
	tkn := tools.GenerateNewShortcutToken()
	a.store[tkn] = user.Id
	return tkn
}
