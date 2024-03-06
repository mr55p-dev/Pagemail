package auth

import (
	"crypto/sha256"

	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

// Validation logic
func (*Authorizer) ValUserAgainstPage(user *db.User, page *db.Page) bool {
	return user.Id == page.UserId
}

func (a *Authorizer) ValCredentialsAgainstUser(email, password string, user *db.User) (isValid bool) {
	emailValid := email == user.Email
	passwordValid := user.Password == a.GenPasswordHash(password)

	return emailValid && passwordValid
}

func (*Authorizer) GenPasswordHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return string(h.Sum(nil))
}

// Session tokens
func (a *Authorizer) GenSessionToken(user *db.User) string {
	tkn := tools.GenerateNewId(50)
	a.store[tkn] = user.Id
	return tkn
}

func (a *Authorizer) ValSessionToken(token string) string {
	uid := a.store[token]
	return uid
}

func (a *Authorizer) RevokeSessionToken(token string) bool {
	_, ok := a.store[token]
	delete(a.store, token)
	return ok
}

func (a *Authorizer) GenShortcutToken(user *db.User) string {
	tkn := tools.GenerateNewShortcutToken()
	a.store[tkn] = user.Id
	return tkn
}
