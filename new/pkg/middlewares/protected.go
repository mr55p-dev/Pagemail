package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/auth"
	"github.com/mr55p-dev/pagemail/pkg/db"
)

func GetProtectedMiddleware(authClient *auth.Authorizer, dbClient *db.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if true {
				user, _ := dbClient.ReadUserByEmail(c.Request().Context(), "ellislunnon@gmail.com")
				c.Set("user", user)
				return next(c)
			}

			cookie, err := c.Cookie(auth.SESS_COOKIE)
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			uid := authClient.CheckToken(cookie.Value)
			if uid == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			requestedId := c.Param("id")
			if requestedId != "" && requestedId != uid {
				return c.NoContent(http.StatusForbidden)
			}

			user, err := dbClient.ReadUserById(c.Request().Context(), uid)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
