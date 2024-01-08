package auth

import (
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/tools"
)

// Get a new token for a given user id
func (a *Authorizer) GetToken(user *db.User) string {
	tkn := tools.GenerateNewId(50)
	TokenStore[tkn] = user.Id
	return tkn
}

// Get a user id from a token
func (a *Authorizer) CheckToken(token string) string {
	uid := TokenStore[token]
	return uid
}

// Revoke a token
func (a *Authorizer) RevokeToken(token string) {
	delete(TokenStore, token)
}
