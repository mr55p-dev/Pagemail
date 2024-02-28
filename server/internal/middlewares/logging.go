package middlewares

import (
	"github.com/labstack/echo/v4"
)

func (p *Provider) GetLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// log the request
		p.log.With(
			"method", c.Request().Method,
			"uri", c.Request().URL.Path,
			"query", c.Request().URL.RawQuery,
		).Info("Request")

		// call the next middleware/handler
		err := next(c)
		if err != nil {
			p.log.With("error", err.Error()).Error("Response")
			return err
		}

		return nil
	}
}
