package middlewares

import (
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
)

func reqWithUser(r *http.Request, user *queries.User) *http.Request {
	userBoundCtx := auth.SetUser(r.Context(), user)
	return r.WithContext(userBoundCtx)
}
