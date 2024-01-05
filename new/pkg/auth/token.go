package auth

import "github.com/mr55p-dev/pagemail/pkg/tools"

// Get a new token for a given user id
func (a *Authorizer) GetToken(id string) string {
	tkn := tools.GenerateNewId(50)
	TokenStore[tkn] = id
	return tkn
}

// Get a user id from a token
func (a *Authorizer) CheckToken(token string) string {
	uid := TokenStore[token]
	return uid
}
