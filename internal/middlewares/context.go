package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
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

func reqWithUser(r *http.Request, user dbqueries.User) *http.Request {
	panic("not implemented")
	// userBoundCtx := db.SetUser(r.Context(), user)
	// return r.WithContext(userBoundCtx)
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

func reqGetUser(r *http.Request) *dbqueries.User {
	panic("not implemented")
}

func reqGetError(r *http.Request) error {
	err, _ := r.Context().Value(userLoadErr).(error)
	return err
}
