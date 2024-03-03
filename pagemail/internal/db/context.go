package db

import "context"

type contextKey string

var userKey contextKey = "user"

func SetUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) *User {
	val, _ := ctx.Value(userKey).(*User)
	return val
}
