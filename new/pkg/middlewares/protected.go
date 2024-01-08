package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/auth"
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/rs/zerolog"
)

func GetProtectedMiddleware(log zerolog.Logger, authClient auth.AbsAuthorizer, dbClient db.AbsClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(auth.SESS_COOKIE)
			if err != nil {
				log.Error().Msg("Could not read session token")
				return c.NoContent(http.StatusUnauthorized)
			}

			uid := authClient.CheckToken(cookie.Value)
			if uid == "" {
				log.Error().Msg("Could not find referenced uid")
				return c.NoContent(http.StatusUnauthorized)
			}

			requestedId := c.Param("id")
			if requestedId != "" && requestedId != uid {
				log.Error().Msgf("Cookie for %s not valid to access user %s", uid, requestedId)
				return c.NoContent(http.StatusForbidden)
			}

			user, err := dbClient.ReadUserById(uid)
			if err != nil {
				log.Error().Msgf("Could not find user with id %s", uid)
				return c.NoContent(http.StatusInternalServerError)
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
