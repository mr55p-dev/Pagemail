package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/render"
)

func (Router) GetRoot(c echo.Context) error {
	var user *db.User
	if u := c.Get("user"); u != nil {
		user = u.(*db.User)
	}
	return render.ReturnRender(c, render.Index(user))
}

func (Router) GetLogin(c echo.Context) error {
	return render.ReturnRender(c, render.Login())
}

func (r *Router) PostLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	user, err := r.DBClient.ReadUserByEmail(c.Request().Context(), email)
	if err != nil {
		return MakeErrorResponse(
			c, http.StatusInternalServerError, err,
		)
	}

	if !r.Authorizer.ValCredentialsAgainstUser(email, password, user) {
		return MakeErrorResponse(c, http.StatusUnauthorized, err)
	}

	sess := r.Authorizer.GenSessionToken(user)

	SetRedirect(c, "/dashboard")
	SetLoginCookie(c, sess)
	return c.NoContent(http.StatusOK)
}

func (Router) GetSignup(c echo.Context) error {
	return render.ReturnRender(c, render.Signup())
}

func (r *Router) PostSignup(c echo.Context) error {
	// Read the form requests
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	passwordRepeat := c.FormValue("password-repeat")
	if password != passwordRepeat {
		return c.String(http.StatusBadRequest, "Passwords do not match")
	}

	// Generate a new user
	user := db.NewUser(email, r.Authorizer.GenPasswordHash(password))
	user.Username = username
	err := r.DBClient.CreateUser(c.Request().Context(), user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Something went wrong")
	}

	// Generate a token for the user from the session manager
	token := r.Authorizer.GenSessionToken(user)

	SetRedirect(c, "/dashboard")
	SetLoginCookie(c, token)
	return c.NoContent(http.StatusOK)
}

func (Router) GetLogout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:   auth.SESS_COOKIE,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	SetRedirect(c, "/")
	return c.NoContent(http.StatusOK)
}

func (r *Router) GetDashboard(c echo.Context) error {
	user := LoadUser(c)
	pages, err := r.DBClient.ReadPagesByUserId(c.Request().Context(), user.Id, 1)
	if err != nil {
		return MakeErrorResponse(c, http.StatusInternalServerError, err)
	}

	return render.ReturnRender(c, render.Dashboard(user, pages))
}

func (r *Router) GetPages(c echo.Context) error {
	user := LoadUser(c)
	page, err := strconv.Atoi(c.QueryParam("p"))
	if err != nil {
		return MakeErrorResponse(c, http.StatusBadRequest, err)
	}

	pages, err := r.DBClient.ReadPagesByUserId(c.Request().Context(), user.Id, page)
	if err != nil {
		return MakeErrorResponse(c, http.StatusInternalServerError, err)
	}
	return render.ReturnRender(c, render.PageList(pages, page))
}

func (r *Router) DeletePages(c echo.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	n, err := r.DBClient.DeletePagesByUserId(c.Request().Context(), user.Id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return render.ReturnRender(c, render.SavePageSuccess(fmt.Sprintf("Deleted %d records", n)))

}

func (r *Router) GetPage(c echo.Context) error {
	user := LoadUser(c)
	pageId := c.Param("page_id")
	page, err := r.DBClient.ReadPage(c.Request().Context(), pageId)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.String(http.StatusInternalServerError, "Failed to get page id")
	}
	if page.UserId != user.Id {
		return c.NoContent(http.StatusForbidden)
	}
	err = r.DBClient.DeletePage(c.Request().Context(), pageId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return render.ReturnRender(c, render.PageCard(page))
}

func (r *Router) PostPage(c echo.Context) error {

	user, ok := c.Get("user").(*db.User)
	if !ok || user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	url := c.FormValue("url")
	if url == "" {
		return render.ReturnRender(c, render.SavePageError("URL field must be present"))
	}

	page := db.NewPage(user.Id, url)
	if err := r.DBClient.CreatePage(c.Request().Context(), page); err != nil {
		return render.ReturnRender(c, render.SavePageError(err.Error()))
	}

	go func(cli *db.Client) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := preview.FetchPreview(ctx, page)
		if err != nil {
			return
		}
		err = cli.UpsertPage(ctx, page)
		if err != nil {
			return
		}
	}(r.DBClient)

	if IsHtmx(c) {
		return render.ReturnRender(c, render.PageCard(page))
	} else {
		return c.String(http.StatusOK, fmt.Sprintf("Added %s succesfully", url))
	}
}

func (r *Router) DeletePage(c echo.Context) error {
	user := LoadUser(c)
	id := c.Param("page_id")
	page, err := r.DBClient.ReadPage(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	if !r.Authorizer.ValUserAgainstPage(user, page) {
		return c.String(http.StatusForbidden, "Permission denied")
	}

	r.DBClient.DeletePage(c.Request().Context(), page.Id)

	return c.NoContent(http.StatusOK)
}

func (r *Router) GetAccountPage(c echo.Context) error {
	user := LoadUser(c)
	return render.ReturnRender(c, render.AccountPage(user))
}

func (r *Router) PutAccount(c echo.Context) error {
	user := LoadUser(c)
	form := new(AccountForm)
	err := c.Bind(form)
	if err != nil {
		return MakeErrorResponse(c, http.StatusBadRequest, err)
	}
	user.Subscribed = form.Subscribed == "on"
	err = r.DBClient.UpdateUser(c.Request().Context(), user)
	if err != nil {
		return MakeErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, "Success!")
}

func (r *Router) GetShortcutToken(c echo.Context) error {
	user := LoadUser(c)
	token := r.Authorizer.GenShortcutToken(user)
	user.ShortcutToken = token
	err := r.DBClient.UpdateUser(c.Request().Context(), user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, token)
}
