package auth

import (
	"crypto/sha256"

	"github.com/mr55p-dev/pagemail/pkg/db"
)

func HashPassword(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return string(h.Sum(nil))
}

func (a *Authorizer) ValidateUser(email, password string, user *db.User) (isValid bool) {
	emailValid := email == user.Email
	passwordValid := user.Password == HashPassword(password)

	return emailValid && passwordValid
}
