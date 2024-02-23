package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
)

func GetShortcutProtected(authClient auth.Authorizer, dbClient *db.Client) echo.MiddlewareFunc {
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

func GetProtectedMiddleware(authClient auth.Authorizer, dbClient *db.Client, block bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(auth.SESS_COOKIE)
			if err != nil || cookie == nil {
				log.ReqErr(c, "Request blocked, no auth token", err)
				if block {
					return c.NoContent(http.StatusUnauthorized)
				} else {
					next(c)
					return nil
				}
			}

			uid := authClient.ValSessionToken(cookie.Value)
			if uid == "" {
				log.ReqError(c, "Request blocked, session not valid")
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
				log.ReqErr(c, "Error, could not read user", err, logging.UserId, uid)
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
