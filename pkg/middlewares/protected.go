package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/auth"
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/logging"
)

func GetProtectedMiddleware(authClient auth.Authorizer, dbClient *db.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(auth.SESS_COOKIE)
			if err != nil || cookie == nil {
				log.ReqErr(c, "Request blocked, no auth token", err)
				return c.NoContent(http.StatusUnauthorized)
			}

			uid := authClient.ValSessionToken(cookie.Value)
			if uid == "" {
				log.ReqError(c, "Request blocked, session not valid")
				return c.NoContent(http.StatusUnauthorized)
			}

			requestedId := c.Param("id")
			if requestedId != "" && requestedId != uid {
				log.ReqError(c, "Request blocked, requested resource does not match session")
				return c.NoContent(http.StatusForbidden)
			}

			user, err := dbClient.ReadUserById(c.Request().Context(), uid)
			if err != nil {
				log.ReqErr(c, "Error, could not read user", err, logging.UserId, uid)
				return c.NoContent(http.StatusInternalServerError)
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
