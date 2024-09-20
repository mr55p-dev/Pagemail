package main

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func PanicError(msg string, err error) {
	if err == nil {
		return
	}
	logger.Error(msg, "error", err.Error())
	panic("Panic generated: " + msg + "; error: " + err.Error())
}

func LogError(msg string, err error) {
	if err == nil {
		return
	}
	logger.Error(msg, "error", err.Error())
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
