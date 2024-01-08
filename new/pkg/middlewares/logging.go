package middlewares

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

var TraceId = "pm-trace-id"

func GetLoggingMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// log the request
			logger.With(
				"method", c.Request().Method,
				"uri", c.Request().URL.Path,
				"query", c.Request().URL.RawQuery,
			).Info("Request")

			// call the next middleware/handler
			err := next(c)
			if err != nil {
				logger.With("error", err.Error()).Error("Response")
				return err
			}

			return nil
		}
	}
}
