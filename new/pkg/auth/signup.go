package auth

import (
	"github.com/mr55p-dev/pagemail/pkg/db"
)

func (a *Authorizer) SignupNewUser(email, password, username string) (token string, err error) {
	user := db.NewUser(email, HashPassword(password))
	user.Username = username
	a.log.Info().Msgf("user: %+v", user)

	err = a.DBClient().InsertUser(user)
	if err != nil {
		return
	}

	token = a.GetToken(user.Id)
	return
}
