package middlewares

import (
	"context"
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
)

type DB interface {
	ReadUserById(context.Context, string) (*db.User, error)
	ReadUserByShortcutToken(context.Context, string) (*db.User, error)
}

type Auth interface {
	ValSessionToken(string) string
}

func GetUserLoader(auth Auth, db DB) MiddlewareFunc {
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

			uid := auth.ValSessionToken(tkn.Value)
			if uid == "" {
				logger.DebugCtx(r.Context(), "Failed to match session cookie with user", "cookie", tkn.Value)
				next.ServeHTTP(w, r)
				return
			}

			user, err := db.ReadUserById(r.Context(), uid)
			if err != nil {
				http.Error(w, "error reading user", http.StatusInternalServerError)
				return
			}

			logger.DebugCtx(r.Context(), "Loaded user from session cookie", "user", user)
			next.ServeHTTP(w, reqWithUser(r, user))
		})
	}
}

func GetShortcutLoader(auth Auth, db DB) MiddlewareFunc {
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

			next.ServeHTTP(w, reqWithUser(r, user))
			return
		})
	}
}

func ProtectRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := reqGetUser(r)
		if user == nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
