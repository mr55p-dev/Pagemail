package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/auth"
	"github.com/mr55p-dev/pagemail/pkg/db"
)

func GetProtectedMiddleware(log *slog.Logger, authClient auth.AbsAuthorizer, dbClient db.AbsClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(auth.SESS_COOKIE)
			if err != nil {
				log.Error("Could not read session token")
				return c.NoContent(http.StatusUnauthorized)
			}

			uid := authClient.CheckToken(cookie.Value)
			if uid == "" {
				log.Error("Could not find referenced uid")
				return c.NoContent(http.StatusUnauthorized)
			}

			requestedId := c.Param("id")
			if requestedId != "" && requestedId != uid {
				log.Error("Cookie for not valid to access user", "cookie_user_id", uid, "request_user_id", requestedId)
				return c.NoContent(http.StatusForbidden)
			}

			user, err := dbClient.ReadUserById(uid)
			if err != nil {
				log.Error("Could not find user with", "cookie_user_id", uid, "request_user_id", requestedId, "user_id", uid)
				return c.NoContent(http.StatusInternalServerError)
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
