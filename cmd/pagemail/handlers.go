package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

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
		err = errors.New("Invalid provider")
	}
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
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

func (h *Handlers) GetApp(c echo.Context) error {
	user, err := h.User(c)
	if err != nil {
		return RenderUserError(c, err)
	}

	return Render(c, http.StatusOK, render.App(*user))
}
