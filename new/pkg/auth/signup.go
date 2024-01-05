package auth

import (
	"github.com/mr55p-dev/pagemail/pkg/db"
)

func (a *Authorizer) SignupNewUser(email, password, username string) error {
	user := db.NewUser(email, password)
	user.Username = username
	return a.DBClient().InsertUser(user)
}
