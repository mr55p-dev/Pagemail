package middlewares

import (
	"net/http"

	"github.com/mr55p-dev/htmx-utils"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
)

type Protector struct {
	auth *auth.Authorizer
	db   *db.Client
	log  *logging.Logger
}

func NewProtector(authorizer *auth.Authorizer, dbclient *db.Client, logger *logging.Logger) *Protector {
	return &Protector{
		auth: authorizer,
		db:   dbclient,
		log:  logger,
	}
}

func (p *Protector) LoadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := p.log.WithRequest(r)
		tkn, err := r.Cookie("pm-auth-tkn")
		if err != nil {
			log.DebugCtx(r.Context(), "No session cookie found")
			next.ServeHTTP(w, reqWithError(r, "could not decode session cookie", http.StatusBadRequest))
			return
		}

		if tkn.Value == "" {
			http.Error(w, "missing session cookie", http.StatusBadRequest)
			return
		}

		uid := p.auth.ValSessionToken(tkn.Value)
		if uid == "" {
			log.DebugCtx(r.Context(), "Failed to match session cookie with user", "cookie", tkn.Value)
			next.ServeHTTP(w, r)
			return
		}

		user, err := p.db.ReadUserById(r.Context(), uid)
		if err != nil {
			http.Error(w, "error reading user", http.StatusInternalServerError)
			return
		}

		log.DebugCtx(r.Context(), "Loaded user from session cookie", "user", user)
		next.ServeHTTP(w, reqWithUser(r, user))
	})
}

func (p *Protector) LoadFromShortcut() hut.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tkn := r.Header.Get("Authorization")
			if tkn == "" {
				http.Error(w, "missing shortcut token", http.StatusBadRequest)
				return
			}

			user, err := p.db.ReadUserByShortcutToken(r.Context(), tkn)
			if err != nil {
				http.Error(w, "invalid shortcut token", http.StatusUnauthorized)
				return
			}

			next(w, reqWithUser(r, user))
			return
		}
	}
}

func (p *Protector) ProtectRoute() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := reqGetUser(r)
			if user == nil {
				err := reqGetError(r)
				p.log.WithRequest(r).WithError(err).ErrorCtx(r.Context(), "User access forbidden")
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
