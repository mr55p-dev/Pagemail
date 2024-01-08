package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/tools"
)

func TraceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(TraceId, tools.GenerateNewId(10))
		return next(c)
	}
}
