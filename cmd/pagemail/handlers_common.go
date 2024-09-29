package main

import (
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/render"
)

const idKey = "id"
const sessionKey = "pm-session"

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

func RenderMsg(ctx echo.Context, statusCode int, msg string) error {
	return Render(ctx, statusCode, render.Message(msg))
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

func GetUser(c echo.Context) queries.User {
	return c.Get("user").(queries.User)
}

// Redirect performs a HXMX-safe redirect
func Redirect(c echo.Context, location string) error {
	c.Response().Header().Set("HX-Location", location)
	if isHTMX(c) {
		return c.NoContent(http.StatusOK)
	}
	return c.Redirect(http.StatusSeeOther, location)
}

// Queries gives access to the db queries object directly
func (h *Handlers) Queries() *queries.Queries {
	return queries.New(h.conn)
}

// User fetches a user object based on the user session cookie.
// Errors when there is no cookie set
func (h *Handlers) User(c echo.Context) (*queries.User, error) {
	sess, err := session.Get(sessionKey, c)
	if err != nil {
		return nil, errors.New("Failed to read user session")
	}
	id, ok := sess.Values[idKey].(string)
	if !ok {
		return nil, errors.New("Invalid id key")
	}
	user, err := h.Queries().ReadUserById(c.Request().Context(), )
	if err != nil {
		return nil, errors.New("Failed to read user")
	}
	return &user, nil
}

// NeedsUser is a middleware func to require a user and set it in context
func (h *Handlers) NeedsUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := h.User(c)
		if err != nil {
			return RenderUserError(c, err)
		}
		c.Set("user", *user)
		return next(c)
	}
}
