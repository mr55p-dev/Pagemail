package auth

import (
	"crypto/subtle"
	"errors"

	"github.com/mr55p-dev/pagemail/internal/tools"
	"golang.org/x/crypto/bcrypt"
)

// Validation logic
func (*Authorizer) ValUserAgainstPage(userID, pageUserID string) bool {
	return userID == pageUserID
}

var (
	ErrInvlaidUsername = errors.New("Incorrect username")
	ErrInvalidPassword = errors.New("Incorrect password")
)

func ValidateUser(userEmail, dbEmail, userPassword, dbPasswordHash []byte) error {
	isValid := subtle.ConstantTimeCompare(userEmail, dbEmail)
	if isValid != 1 {
		return ErrInvlaidUsername
	}

	if err := bcrypt.CompareHashAndPassword(dbPasswordHash, userPassword); err != nil {
		return ErrInvalidPassword
	}

	return nil
}

func HashPassword(pass []byte) []byte {
	pass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		logger.WithError(err).Error("Could not generate password hash")
		panic(err)
	}
	return pass
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
