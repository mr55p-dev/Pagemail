package middlewares

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

func (p *Provider) TraceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := tools.GenerateNewId(10)
		ctx := context.WithValue(c.Request().Context(), "trace-id", id)
		req := c.Request().WithContext(ctx)

		c.SetRequest(req)
		c.Response().Header().Set("X-pm-trace-id", id)
		return next(c)
	}
}
