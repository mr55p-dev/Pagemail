package auth

import (
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/rs/zerolog"
)

type AbsAuthorizer interface {
	DBClient() db.AbsClient
	ValidateUser(email, password string) (isUser bool)
	SignupNewUser(email, password, username string) error
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
