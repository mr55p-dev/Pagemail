package auth

import (
	"context"
	"time"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

func SignupUserIdp(ctx context.Context, queries *dbqueries.Queries, email, name string) (dbqueries.User, error) {
	now := time.Now()
	id := tools.GenerateNewId(10)
	_, shortcutTkn := NewShortcutToken()
	user := dbqueries.User{
		ID:             id,
		Username:       name,
		Email:          email,
		Subscribed:     true,
		ShortcutToken:  shortcutTkn,
		HasReadability: false,
		Created:        now,
		Updated:        now,
	}
	err := queries.CreateUser(ctx, dbqueries.CreateUserParams{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Subscribed:     user.Subscribed,
		ShortcutToken:  user.ShortcutToken,
		HasReadability: user.HasReadability,
		Created:        user.Created,
		Updated:        user.Updated,
	})
	return user, err
}
