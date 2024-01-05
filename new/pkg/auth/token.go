package auth

import "github.com/mr55p-dev/pagemail/pkg/tools"

// Get a new token for a given user id
func (a *Authorizer) GetToken(id string) string {
	tkn := tools.GenerateNewId(50)
	a.log.Debug().Msgf("Created token for user %s (%s)", id, tkn)
	TokenStore[tkn] = id
	return tkn
}

// Get a user id from a token
func (a *Authorizer) CheckToken(token string) string {
	uid, ok := TokenStore[token]
	if !ok {
		a.log.Debug().Msgf("Could not read token for token %s", token)
	} else {
		a.log.Debug().Msgf("Found uid %s for token %s", uid, token)
	}

	return uid
}
