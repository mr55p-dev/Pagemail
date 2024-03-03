package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/mr55p-dev/go-httpit"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
)

type Provider struct {
}

type Protector struct {
	auth auth.Authorizer
	db   *db.Client
	log  *logging.Logger
}

func NewProtector(authorizer auth.Authorizer, dbclient *db.Client, logger *logging.Logger) *Protector {
	return &Protector{
		auth: authorizer,
		db:   dbclient,
		log:  logger,
	}
}

type loadKey string

var userLoadErr loadKey = "user-error"

func requestWithUser(r *http.Request, user *db.User) *http.Request {
	userBoundCtx := db.SetUser(r.Context(), user)
	return r.WithContext(userBoundCtx)
}

func requestWithError(r *http.Request, msg string, code int) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userLoadErr, &httpit.HttpError{
		Msg:  msg,
		Code: code,
	}))
}

func getError(r *http.Request) error {
	err, _ := r.Context().Value(userLoadErr).(error)
	return err
}

func (p *Protector) LoadUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tkn := r.Header.Get("Authorization")
		if tkn == "" {
			next(w, requestWithError(r, "missing user token", http.StatusBadRequest))
			return
		}

		uid := p.auth.ValSessionToken(tkn)
		if uid == "" {
			next(w, requestWithError(r, "token not matched with session", http.StatusBadRequest))
			return
		}

		user, err := p.db.ReadUserById(r.Context(), uid)
		if err != nil {
			http.Error(w, "error reading user", http.StatusInternalServerError)
			return
		}

		next(w, requestWithUser(r, user))
	}
}

func (p *Protector) LoadFromShortcut() httpit.MiddlewareFunc {
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

			next(w, requestWithUser(r, user))
			return
		}
	}
}

func (p *Protector) ProtectRoute() httpit.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := db.GetUser(r.Context())
			if user == nil {
				err := getError(r)
				if err != nil {
					err = errors.New("Missing user")
				}
				httpit.WriteError(w, err)
				return
			}
			next(w, r)
		}
	}
}

// func (p *Provider) GetShortcutProtected(authClient auth.Authorizer, dbClient *db.Client) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			token := c.Request().Header.Get("Authorization")
// 			userId := authClient.ValSessionToken(token)
// 			if userId == "" {
// 				return c.NoContent(http.StatusUnauthorized)
// 			}
// 			user, err := dbClient.ReadUserById(c.Request().Context(), userId)
// 			if err != nil {
// 				return c.NoContent(http.StatusNotFound)
// 			}
// 			c.Set("user", user)
// 			next(c)
// 			return nil
// 		}
// 	}
// }
//
// func (p *Provider) GetProtectedMiddleware(authClient auth.Authorizer, dbClient *db.Client, block bool) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			cookie, err := c.Cookie(auth.SESS_COOKIE)
// 			if err != nil || cookie == nil {
// 				p.log.Errc(c.Request().Context(), "Request blocked, no auth token", err)
// 				if block {
// 					return c.NoContent(http.StatusUnauthorized)
// 				} else {
// 					next(c)
// 					return nil
// 				}
// 			}
//
// 			uid := authClient.ValSessionToken(cookie.Value)
// 			if uid == "" {
// 				p.log.ErrorContext(c.Request().Context(), "Request blocked, session not valid")
// 				if block {
// 					c.Response().Header().Add("HX-Location", "/login")
// 					return c.NoContent(http.StatusUnauthorized)
// 				} else {
// 					next(c)
// 					return nil
// 				}
// 			}
//
// 			user, err := dbClient.ReadUserById(c.Request().Context(), uid)
// 			if err != nil {
// 				p.log.Errc(c.Request().Context(), "Error, could not read user", err)
// 				if block {
// 					return c.NoContent(http.StatusInternalServerError)
// 				} else {
// 					next(c)
// 					return nil
// 				}
// 			}
//
// 			c.Set("user", user)
// 			return next(c)
// 		}
// 	}
// }
