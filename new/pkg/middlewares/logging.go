package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func GetLoggingMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// log the request
			logger.Info().Fields(map[string]interface{}{
				"method": c.Request().Method,
				"uri":    c.Request().URL.Path,
				"query":  c.Request().URL.RawQuery,
			}).Msg("Request")

			// call the next middleware/handler
			err := next(c)
			if err != nil {
				logger.Error().Fields(map[string]interface{}{
					"error": err.Error(),
				}).Msg("Response")
				return err
			}

			return nil
		}
	}
}
