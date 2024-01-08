package middlewares

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

func GetLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// log the request
		slog.With(
			"method", c.Request().Method,
			"uri", c.Request().URL.Path,
			"query", c.Request().URL.RawQuery,
		).Info("Request")

		// call the next middleware/handler
		err := next(c)
		if err != nil {
			slog.With("error", err.Error()).Error("Response")
			return err
		}

		return nil
	}
}
