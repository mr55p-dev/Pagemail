package auth

import "crypto/sha256"

func HashPassword(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return string(h.Sum(nil))
}

func (a *Authorizer) ValidateUser(email, password string) (isUser bool) {
	a.log.Info().Msgf("Requested login by %s", email)
	user, err := a.DBClient().GetUserByEmail(email)
	if err != nil {
		a.log.Error().Msg(err.Error())
		return false
	}

	emailValid := email == user.Email
	passwordValid := user.Password == HashPassword(password)
	a.log.Info().Msg(password)
	a.log.Info().Msg(HashPassword(password))
	a.log.Info().Msg(user.Password)
	if !emailValid {
		a.log.Info().Msgf("Invalid email %s", email)
	}
	if !passwordValid {
		a.log.Info().Msgf("Invalid password %s", password)
	}

	return emailValid && passwordValid
}
