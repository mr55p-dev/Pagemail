package auth

import (
	"crypto/subtle"
	"errors"

	"github.com/gorilla/sessions"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"golang.org/x/crypto/bcrypt"
)

type uidKeyType string

var (
	logger = logging.NewLogger("auth")

	uid        uidKeyType = "uid"
	SessionKey string     = "pm-auth-tkn"

	ErrNoUserInSession = errors.New("No user found in session")
	ErrInvlaidUsername = errors.New("Incorrect username")
	ErrInvalidPassword = errors.New("Incorrect password")
)

func GetId(sess *sessions.Session) string {
	return sess.Values[uid].(string)
}

func SetId(sess *sessions.Session, id string) {
	sess.Values[uid] = id
}

func ValUserAgainstPage(userID, pageUserID string) bool {
	return userID == pageUserID
}

func ValidateUser(userEmail, dbEmail []byte) error {
	isValid := subtle.ConstantTimeCompare(userEmail, dbEmail)
	if isValid != 1 {
		return ErrInvlaidUsername
	}
	return nil
}

func ValidatePassword(userPassword, dbPasswordHash []byte) error {
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
