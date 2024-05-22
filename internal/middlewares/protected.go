package middlewares

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/logging"
)

func GetUserLoader(store sessions.Store, db *sql.DB) MiddlewareFunc {
	logger := logging.NewLogger("middleware-load-user")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := store.Get(r, auth.SessionKey)
			if err != nil {
				logger.WithError(err).DebugCtx(r.Context(), "Failed to load session")
				next.ServeHTTP(w, r)
				return
			}
			uid := auth.GetId(sess)
			user, err := queries.New(db).ReadUserById(r.Context(), uid)
			if err != nil {
				logger.WithError(err).DebugCtx(r.Context(), "Failed to match session cookie with user")
				next.ServeHTTP(w, r)
				return
			}

			logger.DebugCtx(r.Context(), "Loaded user from session cookie")
			next.ServeHTTP(w, reqWithUser(r, &user))
		})
	}
}

func GetShortcutLoader(authorizer sessions.Store, db *sql.DB) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tkn := r.Header.Get("Authorization")
			if tkn == "" {
				http.Error(w, "missing shortcut token", http.StatusBadRequest)
				return
			}

			tokenHash := auth.HashValue([]byte(tkn))
			user, err := queries.New(db).ReadUserByShortcut(r.Context(), tokenHash)
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
		user := auth.GetUser(r.Context())
		if user == nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
