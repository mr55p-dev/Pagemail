package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/logging"
)

func GetLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	log := logging.GetLogger("requestLogger")
	return func(c echo.Context) error {
		// log the request
		log.With(
			"method", c.Request().Method,
			"uri", c.Request().URL.Path,
			"query", c.Request().URL.RawQuery,
		).Info("Request")

		// call the next middleware/handler
		err := next(c)
		if err != nil {
			log.With("error", err.Error()).Error("Response")
			return err
		}

		return nil
	}
}
