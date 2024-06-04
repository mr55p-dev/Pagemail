package auth

import (
	"context"

	"github.com/mr55p-dev/pagemail/db/queries"
)

type contextKey string

var userKey contextKey = "user"
var authenticatedKey = "authenticated"

func SetUser(ctx context.Context, user *queries.User) context.Context {
	ctx = context.WithValue(ctx, userKey, user)
	ctx = context.WithValue(ctx, authenticatedKey, true)
	return ctx
}

func IsAuthenticated(ctx context.Context) bool {
	val, ok := ctx.Value(authenticatedKey).(bool)
	if !ok {
		return false
	}
	return val
}

func GetUser(ctx context.Context) *queries.User {
	val, _ := ctx.Value(userKey).(*queries.User)
	return val
}
