package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/db"
)

type loadKey string

var userLoadErr loadKey = "user-error"

type RequestError struct {
	Message string
	Status  int
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("Request error (%d): %s", r.Status, r.Message)
}

func reqWithUser(r *http.Request, user *db.User) *http.Request {
	userBoundCtx := db.SetUser(r.Context(), user)
	return r.WithContext(userBoundCtx)
}

func reqWithError(r *http.Request, msg string, code int) *http.Request {
	return r.WithContext(
		context.WithValue(
			r.Context(),
			userLoadErr,
			&RequestError{
				Message: msg,
				Status:  code,
			},
		),
	)
}

func reqGetUser(r *http.Request) *db.User {
	return db.GetUser(r.Context())
}

func reqGetError(r *http.Request) error {
	err, _ := r.Context().Value(userLoadErr).(error)
	return err
}
