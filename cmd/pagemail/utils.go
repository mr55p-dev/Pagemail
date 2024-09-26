package main

import (
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

func LogHandlerError(c echo.Context, msg string, err error) {
	if err == nil {
		return
	}
	logger.InfoContext(c.Request().Context(), msg, "error", err.Error())
}
