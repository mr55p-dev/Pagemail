package auth

import (
	"log/slog"

	"github.com/mr55p-dev/pagemail/pkg/db"
)

// maps tokens to user ids
var TokenStore map[string]string

type AbsAuthorizer interface {
	DBClient() db.Client
	ValidateUser(email, password string, user *db.User) (isUser bool)
	SignupNewUser(email, password, username string) (string, error)
	GetToken(string) string
	CheckToken(token string) string
	RevokeToken(token string)
}

type Authorizer struct {
	client db.Client
	log    *slog.Logger
}

func NewAuthorizer(client db.Client, logger *slog.Logger) AbsAuthorizer {
	return &Authorizer{client, logger}
}

func (a *Authorizer) DBClient() db.Client {
	return a.client
}

func init() {
	TokenStore = make(map[string]string)
}
