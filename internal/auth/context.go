package auth

import (
	"context"

	"github.com/mr55p-dev/pagemail/db/queries"
)

type contextKey string

var UserKey contextKey = "user"

func SetUser(ctx context.Context, user *queries.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func GetUser(ctx context.Context) *queries.User {
	val, _ := ctx.Value(UserKey).(*queries.User)
	return val
}
