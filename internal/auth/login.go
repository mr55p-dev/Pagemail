package auth

import (
	"context"
	"errors"

	"github.com/mr55p-dev/pagemail/db/queries"
	"golang.org/x/crypto/bcrypt"
)

func findAuthPlatform(platforms []queries.Auth, platform string) *queries.Auth {
	for _, a := range platforms {
		if a.Platform == "pagemail" {
			out := a
			return &out
		}
	}
	return nil
}

type LoginNativeParams struct {
	Email    string
	Password []byte
}

func LoginNative(ctx context.Context, queries *queries.Queries, params *LoginNativeParams) (*queries.User, error) {
	user, err := queries.ReadUserByEmail(ctx, params.Email)
	if err != nil {
		return nil, errors.New("Invalid email")
	}

	authMethods, err := queries.ReadAuthMethods(ctx, user.ID)
	if err != nil {
		return nil, errors.New("Failed to read auth providers")
	}

	auth := findAuthPlatform(authMethods, "pagemail")
	if auth == nil {
		return nil, errors.New("No auth information found for this identity provider")
	}

	if err := bcrypt.CompareHashAndPassword(auth.PasswordHash, params.Password); err != nil {
		return nil, errors.New("Invalid password")
	}

	return &user, nil
}
