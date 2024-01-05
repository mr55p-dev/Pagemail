package auth

import (
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/rs/zerolog"
)

// maps tokens to user ids
var TokenStore map[string]string

type AbsAuthorizer interface {
	DBClient() db.AbsClient
	ValidateUser(email, password string, user *db.User) (isUser bool)
	SignupNewUser(email, password, username string) (string, error)
	GetToken(string) string
	CheckToken(token string) string
	RevokeToken(token string)
}

type Authorizer struct {
	client db.AbsClient
	log    zerolog.Logger
}

func NewAuthorizer(client db.AbsClient, logger zerolog.Logger) AbsAuthorizer {
	return &Authorizer{client, logger}
}

func (a *Authorizer) DBClient() db.AbsClient {
	return a.client
}

func init() {
	TokenStore = make(map[string]string)
}
