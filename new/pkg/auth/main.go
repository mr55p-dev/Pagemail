package auth

import (
	"log/slog"

	"github.com/mr55p-dev/pagemail/pkg/db"
)

// maps tokens to user ids
var TokenStore map[string]string

type AbsAuthorizer interface {
	ValidateUser(email, password string, user *db.User) (isUser bool)
	GetToken(*db.User) string
	CheckToken(token string) string
	RevokeToken(token string)
}

type Authorizer struct {
	log *slog.Logger
}

func NewAuthorizer() AbsAuthorizer {
	return &Authorizer{}
}

func init() {
	TokenStore = make(map[string]string)
}
