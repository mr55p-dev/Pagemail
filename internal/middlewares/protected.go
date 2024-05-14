package middlewares

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/logging"
)

type DB interface {
	ReadUserById(context.Context, string) (dbqueries.User, error)
	ReadUserByShortcutToken(context.Context, string) (dbqueries.User, error)
}

func GetUserLoader(store sessions.Store, db DB) MiddlewareFunc {
	logger := logging.NewLogger("middleware-load-user")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.WithRequest(r)
			tkn, err := r.Cookie("pm-auth-tkn")
			if err != nil {
				logger.DebugCtx(r.Context(), "No session cookie found")
				next.ServeHTTP(w, r)
				return
			}

			if tkn.Value == "" {
				http.Error(w, "missing session cookie", http.StatusBadRequest)
				return
			}

			sess, _ := store.Get(r, auth.SessionKey)
			uid := auth.GetId(sess)
			user, err := db.ReadUserById(r.Context(), uid)
			if err != nil {
				logger.WithError(err).DebugCtx(r.Context(), "Failed to match session cookie with user", "cookie", tkn.Value)
				next.ServeHTTP(w, r)
				return
			}

			logger.DebugCtx(r.Context(), "Loaded user from session cookie", "user", user)
			next.ServeHTTP(w, reqWithUser(r, &user))
		})
	}
}

func GetShortcutLoader(auth sessions.Store, db DB) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tkn := r.Header.Get("Authorization")
			if tkn == "" {
				http.Error(w, "missing shortcut token", http.StatusBadRequest)
				return
			}

			user, err := db.ReadUserByShortcutToken(r.Context(), tkn)
			if err != nil {
				http.Error(w, "invalid shortcut token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, reqWithUser(r, &user))
			return
		})
	}
}

func ProtectRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := dbqueries.GetUser(r.Context())
		if user == nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
