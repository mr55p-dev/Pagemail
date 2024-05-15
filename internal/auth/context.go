package auth

import (
	"context"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
)

type contextKey string

var UserKey contextKey = "user"

func SetUser(ctx context.Context, user *dbqueries.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func GetUser(ctx context.Context) *dbqueries.User {
	val, _ := ctx.Value(UserKey).(*dbqueries.User)
	return val
}
