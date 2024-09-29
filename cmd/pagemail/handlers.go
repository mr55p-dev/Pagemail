package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jordan-wright/email"
	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/cmd/pagemail/urls"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/render"
)

// Handlers wraps all the route handlers
type Handlers struct {
	conn  *pgx.Conn
	store *sessions.CookieStore
	mail  *email.Pool
}

type PageComponent func(user *queries.User) templ.Component

// GetPage will render a simple page
func (h *Handlers) GetPage(component PageComponent) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, _ := h.User(c)
		return Render(c, http.StatusOK, component(user))
	}
}

func (h *Handlers) GetLogout(c echo.Context) error {
	sess, err := h.store.Get(c.Request(), sessionKey)
	if err != nil {
		LogError("Failed to get user session", err)
	}
	sess.Options.MaxAge = -1
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		LogError("Failed to write session", err)
	}
	return Redirect(c, "/")
}

func (h *Handlers) PostLogin(c echo.Context) error {
	provider := c.QueryParam("provider")
	err := c.Request().ParseForm()
	if err != nil {
		msg := "Failed to parse form data"
		LogHandlerError(c, msg, err)
		return RenderError(c, http.StatusBadRequest, msg)
	}

	var user queries.User
	switch provider {
	case "native":
		user, err = h.Queries().ReadUserWithCredential(c.Request().Context(), queries.ReadUserWithCredentialParams{
			Email:    c.FormValue("email"),
			Platform: "pagemail",
			Crypt:    c.FormValue("password"),
		})
	default:
		err = errors.New("Invalid provider")
	}
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}
	if user.ID == uuid.Nil {
		return RenderError(c, http.StatusBadRequest, "Invalid username or password")
	}

	if err := h.CreateSession(c, &user); err != nil {
		return err
	}

	return Redirect(c, urls.App)
}

func (h *Handlers) PostSignup(c echo.Context) error {
	provider := c.QueryParam("provider")
	err := c.Request().ParseForm()
	if err != nil {
		msg := "Failed to parse form data"
		LogHandlerError(c, msg, err)
		return RenderError(c, http.StatusBadRequest, msg)
	}

	var user queries.User
	switch provider {
	case "native":
		err := h.WrapTx(c, func(ctx context.Context, q *queries.Queries) error {
			user, err = q.CreateUser(ctx, queries.CreateUserParams{
				Email:    c.FormValue("email"),
				Username: c.FormValue("username"),
			})
			if err != nil {
				return fmt.Errorf("Failed to create user: %w", err)
			}
			err = q.CreateLocalAuth(ctx, queries.CreateLocalAuthParams{
				UserID: user.ID,
				Crypt:  c.FormValue("password"),
			})
			return nil
		})
		if err != nil {
			LogHandlerError(c, "Failed to create user", err)
			return RenderGenericError(c)
		}
	default:
		err = errors.New("Invalid provider")
	}

	if err := h.CreateSession(c, &user); err != nil {
		return err
	}

	return Redirect(c, urls.App)
}

func (h *Handlers) PostPage(c echo.Context) error {
	user := GetUser(c)

	// load the form url
	err := c.Request().ParseForm()
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed to parse form")
	}
	url, err := url.Parse(c.FormValue("url"))
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "The provided URL is not valid")
	}
	pageArg := queries.CreatePageWithPreviewParams{
		ID:          tools.NewPageId(),
		UserID:      user.ID,
		Url:         c.FormValue("url"),
		Title:       pgtype.Text{},
		Description: pgtype.Text{},
	}

	// Load the preview
	pageData, err := GetPreview(url)
	if err != nil {
		LogHandlerError(c, "Could not get preview", err)
	} else {
		if pageData.Title != "" {
			pageArg.Title = pgtype.Text{String: pageData.Title, Valid: true}
		}
		if pageData.Description != "" {
			pageArg.Description = pgtype.Text{String: pageData.Description, Valid: true}
		}
	}

	// Create the page
	page, err := h.Queries().CreatePageWithPreview(c.Request().Context(), pageArg)
	if err != nil {
		LogHandlerError(c, "Could not create a page", err)
		return RenderError(c, http.StatusInternalServerError, "Failed to create page")
	}

	return Render(c, http.StatusCreated, render.PageCard(page))
}

func (h *Handlers) DeletePage(c echo.Context) error {
	user := GetUser(c)
	pageId := c.Param("id")
	// parse the page id hex into bytes
	cnt, err := h.Queries().DeletePageForUser(c.Request().Context(), queries.DeletePageForUserParams{
		ID:     pageId,
		UserID: user.ID,
	})
	if err != nil {
		LogError("Failed to delete page", err)
		return RenderGenericError(c)
	} else if cnt == 0 {
		return RenderError(c, http.StatusNotFound, "No such page found")
	} else {
		return c.NoContent(http.StatusOK)
	}
}

func (h *Handlers) GetApp(c echo.Context) error {
	user := GetUser(c)
	pages, err := h.Queries().ReadPagesByUserId(c.Request().Context(), queries.ReadPagesByUserIdParams{
		UserID: user.ID,
		Limit:  30,
		Offset: 0,
	})
	if err != nil {
		return RenderError(c, http.StatusOK, "Failed to read pages")
	}
	return Render(c, http.StatusOK, render.App(&user, pages))
}
