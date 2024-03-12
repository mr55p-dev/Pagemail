package db

import "context"

type contextKey string

var UserKey contextKey = "user"

func SetUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func GetUser(ctx context.Context) *User {
	val, _ := ctx.Value(UserKey).(*User)
	return val
}
