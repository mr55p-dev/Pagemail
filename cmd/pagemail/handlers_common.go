package main

import (
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/render"
)

// Render displays a templ component
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

// RenderError displays an error
func RenderError(ctx echo.Context, statusCode int, err string) error {
	return Render(ctx, statusCode, render.Error(err))
}

// RenderGenericError displays a basic 500 error
func RenderGenericError(ctx echo.Context) error {
	return RenderError(ctx, http.StatusInternalServerError, "Something went wrong.")
}

// RenderUserError displays an error when the user is not found
func RenderUserError(ctx echo.Context, err error) error {
	return RenderError(ctx, http.StatusBadRequest, err.Error())
}

// isHTMX checks if the request object is based on htmx
func isHTMX(c echo.Context) bool {
	return c.Request().Header.Get("Hx-Request") == "true"
}

// Redirect performs a HXMX-safe redirect
func Redirect(c echo.Context, location string) error {
	c.Response().Header().Set("HX-Location", location)
	var status = http.StatusSeeOther
	if isHTMX(c) {
		status = http.StatusOK
	}
	return c.Redirect(status, location)
}

// Queries gives access to the db queries object directly
func (h *Handlers) Queries() *queries.Queries {
	return queries.New(h.conn)
}

// User fetches a user object based on the user session cookie.
// Errors when there is no cookie set
func (h *Handlers) User(c echo.Context) (*queries.User, error) {
	sess, err := session.Get("pm-session", c)
	if err != nil {
		return nil, errors.New("Failed to read user session")
	}
	id := sess.Values["id"].(string)
	user, err := h.Queries().ReadUserById(c.Request().Context(), id)
	if err != nil {
		return nil, errors.New("Failed to read user")
	}
	return &user, nil
}
