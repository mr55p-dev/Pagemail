package auth

import "github.com/mr55p-dev/pagemail/internal/db"

func (*TestAuthorizer) ValCredentialsAgainstUser(email, password string, user *db.User) bool {
	return true
}

func (*TestAuthorizer) ValUserAgainstPage(user *db.User, page *db.Page) bool {
	return true
}

func (*TestAuthorizer) GenPasswordHash(pass string) string {
	return pass
}

func (a *TestAuthorizer) GenSessionToken(user *db.User) string {
	return "PM_SESSION_TOKEN"
}

func (a *TestAuthorizer) ValSessionToken(token string) (userId string) {
	return "123"
}

func (a *TestAuthorizer) RevokeSessionToken(token string) (ok bool) {
	return true
}

func (a *TestAuthorizer) GenShortcutToken(*db.User) string {
	return "123"
}
