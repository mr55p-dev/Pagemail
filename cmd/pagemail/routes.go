package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/mr55p-dev/htmx-utils"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/preview"
	"github.com/mr55p-dev/pagemail/internal/render"
)

func staticRender(component templ.Component, w http.ResponseWriter, r *http.Request) {
	err := component.Render(r.Context(), w)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Error rendering response")
		http.Error(w, "Error rendering response", http.StatusInternalServerError)
	}
}

func requestBind[T any](w http.ResponseWriter, r *http.Request) *T {
	out := new(T)
	err := Bind(out, r)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to bind request")
		http.Error(w, "Failed to bind request", http.StatusBadRequest)
		return nil
	}
	return out
}

func serverError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func GetLoginCookie(val string) *http.Cookie {
	return &http.Cookie{
		Name:     auth.SESS_COOKIE,
		Value:    val,
		Path:     "/",
		MaxAge:   864000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (Router) GetRoot(w http.ResponseWriter, r *http.Request) {
	user := db.GetUser(r.Context())
	err := render.Index(user).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render index", http.StatusInternalServerError)
		return
	}
}

func (Router) GetLogin(w http.ResponseWriter, r *http.Request) {
	staticRender(render.Login(), w, r)
}

type PostLoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (router *Router) PostLogin(w http.ResponseWriter, r *http.Request) {
	req := requestBind[PostLoginRequest](w, r)
	if req == nil {
		return
	}

	logger.DebugCtx(r.Context(), "Received bound data", "email", req.Email, "req", req)
	user, err := router.DBClient.ReadUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !router.Authorizer.ValCredentialsAgainstUser(req.Email, req.Password, user) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sess := router.Authorizer.GenSessionToken(user)
	cookie := GetLoginCookie(sess)

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	return
}

func (Router) GetSignup(w http.ResponseWriter, r *http.Request) {
	staticRender(render.Signup(), w, r)
	return
}

type PostSignupRequest struct {
	Username       string `form:"username"`
	Email          string `form:"email"`
	Password       string `form:"password"`
	PasswordRepeat string `form:"password-repeat"`
}

func (router *Router) PostSignup(w http.ResponseWriter, r *http.Request) {
	req := requestBind[PostSignupRequest](w, r)
	if req == nil {
		return
	}

	// Read the form requests
	if req.Password != req.PasswordRepeat {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Generate a new user
	user := db.NewUser(req.Email, router.Authorizer.GenPasswordHash(req.Password))
	user.Username = req.Username
	err := router.DBClient.CreateUser(r.Context(), user)
	if err != nil {
		serverError(w, r)
		return
	}

	// Generate a token for the user from the session manager
	token := router.Authorizer.GenSessionToken(user)
	cookie := GetLoginCookie(token)

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (Router) GetLogout(w http.ResponseWriter, r *http.Request) hut.Writer {
	cookie := http.Cookie{
		Name:   auth.SESS_COOKIE,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	return hut.Redirect("/")
}

func (router *Router) GetDashboard(w http.ResponseWriter, r *http.Request) hut.Writer {
	user := db.GetUser(r.Context())
	pages, err := router.DBClient.ReadPagesByUserId(r.Context(), user.Id, 1)
	if err != nil {
		return hut.Error(err, http.StatusInternalServerError)
	}

	return hut.Component(render.Dashboard(user, pages))
}

type GetPagesRequest struct {
	Page string `query:"p"`
}

func (router *Router) GetPages(w http.ResponseWriter, r *http.Request, req *GetPagesRequest) hut.Writer {
	user := db.GetUser(r.Context())
	page, err := strconv.Atoi(req.Page)
	if err != nil {
		return hut.Error(err, http.StatusBadRequest)
	}

	pages, err := router.DBClient.ReadPagesByUserId(r.Context(), user.Id, page)
	if err != nil {
		return hut.Error(err, http.StatusInternalServerError)
	}
	return hut.Component(render.PageList(pages, page))
}

func (router *Router) DeletePages(w http.ResponseWriter, r *http.Request) hut.Writer {
	user := db.GetUser(r.Context())

	n, err := router.DBClient.DeletePagesByUserId(r.Context(), user.Id)
	if err != nil {
		return hut.Error(err, http.StatusInternalServerError)
	}
	return hut.Component(render.SavePageSuccess(fmt.Sprintf("Deleted %d pages", n)))

}

type GetPageRequest struct {
	PageId string `param:"page_id"`
}

func (router *Router) GetPage(w http.ResponseWriter, r *http.Request, req *GetPageRequest) hut.Writer {
	user := db.GetUser(r.Context())
	page, err := router.DBClient.ReadPage(r.Context(), req.PageId)
	if err != nil {
		return hut.ErrorMsg("Failed to get page id", http.StatusInternalServerError)
	}
	if page.UserId != user.Id {
		return hut.ErrorMsg("Not found", http.StatusNotFound)
	}
	err = router.DBClient.DeletePage(r.Context(), req.PageId)
	if err != nil {
		return hut.Error(err, http.StatusInternalServerError)
	}
	return hut.Component(render.PageCard(page))
}

type PostPageRequest struct {
	Url string `form:"url"`
}

func (router *Router) PostPage(w http.ResponseWriter, r *http.Request, req *PostPageRequest) hut.Writer {
	user := db.GetUser(r.Context())
	url := req.Url
	if url == "" {
		return hut.Component(render.SavePageError("URL field must be present"))
	}

	page := db.NewPage(user.Id, url)
	if err := router.DBClient.CreatePage(r.Context(), page); err != nil {
		return hut.Component(render.SavePageError(err.Error()))
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
	}(router.DBClient)

	return hut.Component(render.PageCard(page))
}

type DeletePageRequest struct {
	PageId string `query:"page_id"`
}

func (router *Router) DeletePage(w http.ResponseWriter, r *http.Request, req *DeletePageRequest) hut.Writer {
	user := db.GetUser(r.Context())
	page, err := router.DBClient.ReadPage(r.Context(), req.PageId)
	if err != nil {
		return hut.ErrorMsg("Something went wrong", http.StatusInternalServerError)
	}

	if !router.Authorizer.ValUserAgainstPage(user, page) {
		return hut.ErrorMsg("Permission denied", http.StatusForbidden)
	}

	router.DBClient.DeletePage(r.Context(), page.Id)

	return nil
}

func (router *Router) GetAccountPage(w http.ResponseWriter, r *http.Request) hut.Writer {
	user := db.GetUser(r.Context())
	return hut.Component(render.AccountPage(user))
}

func (router *Router) PutAccount(w http.ResponseWriter, r *http.Request) hut.Writer {
	user := db.GetUser(r.Context())
	form := new(AccountData)
	user.Subscribed = form.Subscribed == "on"
	err := router.DBClient.UpdateUser(r.Context(), user)
	if err != nil {
		return hut.Error(err, http.StatusInternalServerError)
	}

	return nil
}

func (router *Router) GetShortcutToken(w http.ResponseWriter, r *http.Request) hut.Writer {
	user := db.GetUser(r.Context())
	token := router.Authorizer.GenShortcutToken(user)
	user.ShortcutToken = token
	err := router.DBClient.UpdateUser(r.Context(), user)
	if err != nil {
		return hut.Error(err, http.StatusInternalServerError)
	}
	return hut.String(token, http.StatusOK)
}
