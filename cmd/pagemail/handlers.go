package main

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/jordan-wright/email"
	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/cmd/pagemail/urls"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/render"
)

// Handlers wraps all the route handlers
type Handlers struct {
	conn  *sql.DB
	store *sessions.CookieStore
	mail  *email.Pool
}

func (*Handlers) GetIndex(c echo.Context) error {
	return Render(c, http.StatusOK, render.Index())
}

func (*Handlers) GetLogin(c echo.Context) error {
	return Render(c, http.StatusOK, render.Login())
}

func (*Handlers) GetSignup(c echo.Context) error {
	return Render(c, http.StatusOK, render.Signup())
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

	sess, err := h.store.Get(c.Request(), sessionKey)
	sess.Values[idKey] = user.ID
	if err != nil {
		LogHandlerError(c, "Failed to get user session", err)
		return RenderGenericError(c)
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		LogHandlerError(c, "Failed to save user session", err)
		return RenderGenericError(c)
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
		ID:           tools.NewPageId(),
		UserID:       user.ID,
		Url:          c.FormValue("url"),
		Title:        sql.NullString{},
		Description:  sql.NullString{},
		PreviewState: PREVIEW_UNKNOWN,
	}

	// Load the preview
	pageData, err := GetPreview(url)
	if err != nil {
		LogHandlerError(c, "Could not get preview", err)
		pageArg.PreviewState = PREVIEW_FAILURE
	} else {
		pageArg.PreviewState = PREVIEW_SUCCESS
		if pageData.Title != "" {
			pageArg.Title = sql.NullString{String: pageData.Title, Valid: true}
		}
		if pageData.Description != "" {
			pageArg.Description = sql.NullString{String: pageData.Description, Valid: true}
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
	return Render(c, http.StatusOK, render.App(user, pages))
}
