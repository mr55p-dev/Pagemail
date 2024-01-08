package l

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

var TraceId = "pm-trace-id"

func WithTrace(log *slog.Logger, c echo.Context) *slog.Logger {
	return log.With(TraceId, c.Get(TraceId))
}
