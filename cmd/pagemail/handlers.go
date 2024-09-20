package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

func RenderUserError(ctx echo.Context, err error) error {
	return RenderError(ctx, http.StatusBadRequest, err.Error())
}

func isHTMX(c echo.Context) bool {
	return c.Request().Header.Get("Hx-Request") == "true"
}

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

	sess, err := h.store.Get(c.Request(), "pm-session")
	sess.Values["id"] = user.ID
	if err != nil {
		LogHandlerError(c, "Failed to get user session", err)
		return RenderGenericError(c)
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		LogHandlerError(c, "Failed to save user session", err)
		return RenderGenericError(c)
	}

	return Redirect(c, "/app")
}

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

func (h *Handlers) GetApp(c echo.Context) error {
	user, err := h.User(c)
	if err != nil {
		return RenderUserError(c, err)
	}

	return Render(c, http.StatusOK, render.App(*user))
}
