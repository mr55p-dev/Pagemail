package middlewares

import (
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
)

func reqWithUser(r *http.Request, user *dbqueries.User) *http.Request {
	userBoundCtx := auth.SetUser(r.Context(), user)
	return r.WithContext(userBoundCtx)
}
