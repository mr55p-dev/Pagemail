package auth

import (
	"crypto/subtle"
	"encoding/gob"
	"errors"

	"github.com/gorilla/sessions"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"golang.org/x/crypto/bcrypt"
)

type uidKeyType string

var (
	logger = logging.NewLogger("auth")

	uid        uidKeyType = "uid"
	SessionKey string     = "SessionCookie"

	ErrNoUserInSession = errors.New("No user found in session")
	ErrInvlaidUsername = errors.New("Incorrect username")
	ErrInvalidPassword = errors.New("Incorrect password")
)

func init() {
	gob.Register(uid)
	gob.Register(G_Token{})
}

func GetId(sess *sessions.Session) string {
	if sess == nil {
		return ""
	}
	val, ok := sess.Values[uid].(string)
	if !ok {
		return ""
	}
	return val
}

func SetId(sess *sessions.Session, id string) {
	sess.Values[uid] = id
}

func DelId(sess *sessions.Session) {
	delete(sess.Values, uid)
}

func ValidateEmail(userEmail, dbEmail []byte) bool {
	isValid := subtle.ConstantTimeCompare(userEmail, dbEmail)
	if isValid != 1 {
		return false
	}
	return true
}

func CheckSubmittedPasswords(pass, passRepeat []byte) error {
	if ok := subtle.ConstantTimeCompare(pass, passRepeat); ok != 1 {
		return pmerror.ErrDiffPasswords
	}
	return nil
}

func ValidatePassword(userPassword, dbPasswordHash []byte) bool {
	if err := bcrypt.CompareHashAndPassword(dbPasswordHash, userPassword); err != nil {
		return false
	}
	return true
}

func HashPassword(pass []byte) []byte {
	pass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		logger.WithError(err).Error("Could not generate password hash")
		panic(err)
	}
	return pass
}
