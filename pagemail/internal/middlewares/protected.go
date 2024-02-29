package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
)

type Provider struct {
	log *logging.Logger
}

func New(log *logging.Logger) *Provider {
	return &Provider{
		log: log,
	}
}

func (p *Provider) GetShortcutProtected(authClient auth.Authorizer, dbClient *db.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			userId := authClient.ValSessionToken(token)
			if userId == "" {
				return c.NoContent(http.StatusUnauthorized)
			}
			user, err := dbClient.ReadUserById(c.Request().Context(), userId)
			if err != nil {
				return c.NoContent(http.StatusNotFound)
			}
			c.Set("user", user)
			next(c)
			return nil
		}
	}
}

func (p *Provider) GetProtectedMiddleware(authClient auth.Authorizer, dbClient *db.Client, block bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(auth.SESS_COOKIE)
			if err != nil || cookie == nil {
				p.log.Errc(c.Request().Context(), "Request blocked, no auth token", err)
				if block {
					return c.NoContent(http.StatusUnauthorized)
				} else {
					next(c)
					return nil
				}
			}

			uid := authClient.ValSessionToken(cookie.Value)
			if uid == "" {
				p.log.ErrorContext(c.Request().Context(), "Request blocked, session not valid")
				if block {
					c.Response().Header().Add("HX-Location", "/login")
					return c.NoContent(http.StatusUnauthorized)
				} else {
					next(c)
					return nil
				}
			}

			user, err := dbClient.ReadUserById(c.Request().Context(), uid)
			if err != nil {
				p.log.Errc(c.Request().Context(), "Error, could not read user", err)
				if block {
					return c.NoContent(http.StatusInternalServerError)
				} else {
					next(c)
					return nil
				}
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
