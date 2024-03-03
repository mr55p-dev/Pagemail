package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/mr55p-dev/go-httpit"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/preview"
	"github.com/mr55p-dev/pagemail/internal/render"
)

func (Router) GetRoot(ctx httpit.Context) (templ.Component, error) {
	user := db.GetUser(ctx)
	return render.Index(user), nil
}

func (Router) GetLogin(ctx httpit.Context) (templ.Component, error) {
	return render.Login(), nil
}

type PostLoginRequest struct {
	email    string `form:"email"`
	password string `form:"password"`
}

func (r *Router) PostLogin(ctx httpit.Context, req *PostLoginRequest) error {
	user, err := r.DBClient.ReadUserByEmail(ctx, req.email)

	if !r.Authorizer.ValCredentialsAgainstUser(req.email, req.password, user) {
		return httpit.Error(err.Error(), http.StatusUnauthorized)
	}

	sess := r.Authorizer.GenSessionToken(user)
	cookie := GetLoginCookie(sess)

	ctx.SetRedirect("/dashboard")
	ctx.SetCookie(cookie)
	return nil
}

func (Router) GetSignup(ctx httpit.Context) (templ.Component, error) {
	return render.Signup(), nil
}

type PostSignupRequest struct {
	username       string `form:"username"`
	email          string `form:"email"`
	password       string `form:"password"`
	passwordRepeat string `form:"password-repeat"`
}

func (r *Router) PostSignup(ctx httpit.Context, req *PostSignupRequest) error {
	// Read the form requests
	if req.password != req.passwordRepeat {
		return httpit.Error("Passwords do not match", http.StatusBadRequest)
	}

	// Generate a new user
	user := db.NewUser(req.email, r.Authorizer.GenPasswordHash(req.password))
	user.Username = req.username
	err := r.DBClient.CreateUser(ctx, user)
	if err != nil {
		return httpit.Error("Something went wrong", http.StatusBadRequest)
	}

	// Generate a token for the user from the session manager
	token := r.Authorizer.GenSessionToken(user)
	cookie := GetLoginCookie(token)

	ctx.SetRedirect("/dashboard")
	ctx.SetCookie(cookie)
	return nil
}

func (Router) GetLogout(ctx httpit.Context) error {
	cookie := http.Cookie{
		Name:   auth.SESS_COOKIE,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	ctx.SetCookie(cookie)
	ctx.SetRedirect("/")
	return nil
}

func (r *Router) GetDashboard(ctx httpit.Context) (templ.Component, error) {
	user := db.GetUser(ctx)
	pages, err := r.DBClient.ReadPagesByUserId(ctx, user.Id, 1)
	if err != nil {
		return nil, httpit.Error(err.Error(), http.StatusInternalServerError)
	}

	return render.Dashboard(user, pages), nil
}

type GetPagesRequest struct {
	Page string `query:"p"`
}

func (r *Router) GetPages(ctx httpit.Context, req *GetPagesRequest) (templ.Component, error) {
	user := db.GetUser(ctx)
	page, err := strconv.Atoi(req.Page)
	if err != nil {
		return nil, httpit.Error(err.Error(), http.StatusBadRequest)
	}

	pages, err := r.DBClient.ReadPagesByUserId(ctx, user.Id, page)
	if err != nil {
		return nil, httpit.Error(err.Error(), http.StatusInternalServerError)
	}
	return render.PageList(pages, page), nil
}

func (r *Router) DeletePages(ctx httpit.Context) (templ.Component, error) {
	user := db.GetUser(ctx)

	n, err := r.DBClient.DeletePagesByUserId(ctx, user.Id)
	if err != nil {
		return nil, httpit.Error(err.Error(), http.StatusInternalServerError)
	}
	return render.SavePageSuccess(fmt.Sprintf("Deleted %d pages", n)), nil

}

type GetPageRequest struct {
	PageId string `param:"page_id"`
}

func (r *Router) GetPage(ctx httpit.Context, req *GetPageRequest) (templ.Component, error) {
	user := db.GetUser(ctx)
	page, err := r.DBClient.ReadPage(ctx, req.PageId)
	if err != nil {
		return nil, httpit.Error("Failed to get page id", http.StatusInternalServerError)
	}
	if page.UserId != user.Id {
		return nil, httpit.Error("Not found", http.StatusNotFound)
	}
	err = r.DBClient.DeletePage(ctx, req.PageId)
	if err != nil {
		return nil, httpit.Error(err.Error(), http.StatusInternalServerError)
	}
	return render.PageCard(page), nil
}

type PostPageRequest struct {
	Url string `form:"url"`
}

func (r *Router) PostPage(ctx httpit.Context, req *PostPageRequest) (templ.Component, error) {
	user := db.GetUser(ctx)
	url := req.Url
	if url == "" {
		return render.SavePageError("URL field must be present"), nil
	}

	page := db.NewPage(user.Id, url)
	if err := r.DBClient.CreatePage(ctx, page); err != nil {
		return render.SavePageError(err.Error()), nil
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

	return render.PageCard(page), nil
}

type DeletePageRequest struct {
	PageId string `query:"page_id"`
}

func (r *Router) DeletePage(ctx httpit.Context, req *DeletePageRequest) error {
	user := db.GetUser(ctx)
	page, err := r.DBClient.ReadPage(ctx, req.PageId)
	if err != nil {
		return httpit.Error("Something went wrong", http.StatusInternalServerError)
	}

	if !r.Authorizer.ValUserAgainstPage(user, page) {
		return httpit.Error("Permission denied", http.StatusForbidden)
	}

	r.DBClient.DeletePage(ctx, page.Id)

	return nil
}

func (r *Router) GetAccountPage(ctx httpit.Context) (templ.Component, error) {
	user := db.GetUser(ctx)
	return render.AccountPage(user), nil
}

func (r *Router) PutAccount(ctx httpit.Context) (templ.Component, error) {
	user := db.GetUser(ctx)
	form := new(AccountData)
	user.Subscribed = form.Subscribed == "on"
	err := r.DBClient.UpdateUser(ctx, user)
	if err != nil {
		return nil, httpit.Error(err.Error(), http.StatusInternalServerError)
	}

	return templ.NopComponent, nil
}

func (r *Router) GetShortcutToken(ctx httpit.Context) error {
	user := db.GetUser(ctx)
	token := r.Authorizer.GenShortcutToken(user)
	user.ShortcutToken = token
	err := r.DBClient.UpdateUser(ctx, user)
	if err != nil {
		return httpit.Error(err.Error(), http.StatusInternalServerError)
	}
	// return c.String(http.StatusOK, token)
	return nil
}
