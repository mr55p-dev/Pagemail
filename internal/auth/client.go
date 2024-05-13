package auth

import (
	"crypto/sha256"

	"github.com/mr55p-dev/pagemail/internal/tools"
)

// Validation logic
func (*Authorizer) ValUserAgainstPage(userID, pageUserID string) bool {
	return userID == pageUserID
}

func (a *Authorizer) ValCredentialsAgainstUser(email, password, dbEmail, dbPassword string) (isValid bool) {
	emailValid := email == dbEmail
	passwordValid := string(dbPassword) == a.GenPasswordHash(password)

	return emailValid && passwordValid
}

func (*Authorizer) GenPasswordHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return string(h.Sum(nil))
}

// Session tokens
func (a *Authorizer) GenSessionToken(userID string) string {
	tkn := tools.GenerateNewId(50)
	a.store[tkn] = userID
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

func (a *Authorizer) GenShortcutToken(userID string) string {
	tkn := tools.GenerateNewShortcutToken()
	a.store[tkn] = userID
	return tkn
}
