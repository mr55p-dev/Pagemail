package auth

func HashPassword(pass string) string {
	return pass
}

func (a *Authorizer) ValidateUser(email, password string) (isUser bool) {
	a.log.Info().Msgf("Requested login by %s", email)
	user, err := a.DBClient().GetUserByEmail(email)
	if err != nil {
		return false
	}

	if user.Email == email && user.Password == HashPassword(password) {
		return true
	}
	return
}
