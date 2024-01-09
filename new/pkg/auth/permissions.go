package auth

import "github.com/mr55p-dev/pagemail/pkg/db"

func (a *Authorizer) CheckPagePermission(user *db.User, page *db.Page) bool {
	return user.Id == page.UserId
}
