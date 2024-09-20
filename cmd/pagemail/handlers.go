package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/render"
)

// Handlers wraps all the route handlers
type Handlers struct {
	conn  *sql.DB
	store *sessions.CookieStore
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func RenderError(ctx echo.Context, statusCode int, err string) error {
	return Render(ctx, statusCode, render.Error(err))
}

func RenderGenericError(ctx echo.Context) error {
	return RenderError(ctx, http.StatusInternalServerError, "Something went wrong.")
}

// Queries gives access to the db queries object directly
func (h *Handlers) Queries() *queries.Queries {
	return queries.New(h.conn)
}

// GetIndex constructs the root page
func (s *Handlers) GetIndex(c echo.Context) error {
	return Render(c, http.StatusOK, render.Index())
}

func (s *Handlers) GetLogin(c echo.Context) error {
	return Render(c, http.StatusOK, render.Login())
}

func (h *Handlers) PostLogin(c echo.Context) error {
	provider := c.QueryParam("provider")
	err := c.Request().ParseForm()
	if err != nil {
		msg := "Failed to parse form data"
		LogHandlerError(c, msg, err)
		return RenderError(c, http.StatusBadRequest, msg)
	}

	var user *queries.User
	switch provider {
	case "native":
		user, err = auth.LoginNative(c.Request().Context(), h.Queries(), &auth.LoginNativeParams{
			Email:    c.FormValue("email"),
			Password: []byte(c.FormValue("password")),
		})
	default:
		err = fmt.Errorf("Invalid provider")
	}
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid provider")
	}
	if user == nil {
		LogHandlerError(c, "Could not find user, but no error produced", errors.New("No error"))
		return RenderError(c, http.StatusInternalServerError, "User not found")
	}

	session, err := h.store.Get(c.Request(), user.ID)
	if err != nil {
		LogHandlerError(c, "Failed to get user session", err)
		return RenderGenericError(c)
	}

	if err := session.Save(c.Request(), c.Response()); err != nil {
		LogHandlerError(c, "Failed to save user session", err)
		return RenderGenericError(c)
	}

	return c.Redirect(http.StatusSeeOther, "/app")
}
