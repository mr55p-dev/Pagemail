package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/preview"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/tools"
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

func genericResponse(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
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
	user := dbqueries.GetUser(r.Context())
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

	if !router.Authorizer.ValCredentialsAgainstUser(req.Email, req.Password, user.Email, user.Password.([]byte)) {
		genericResponse(w, http.StatusUnauthorized)
		return
	}

	sess := router.Authorizer.GenSessionToken(user.ID)
	cookie := GetLoginCookie(sess)

	http.SetCookie(w, cookie)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
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
	now := time.Now()
	user := dbqueries.CreateUserParams{
		ID:             tools.GenerateNewId(10),
		Username:       req.Username,
		Email:          req.Email,
		Password:       router.Authorizer.GenPasswordHash(req.Password),
		Avatar:         sql.NullString{},
		Subscribed:     false,
		ShortcutToken:  tools.GenerateNewShortcutToken(),
		HasReadability: false,
		Created:        now,
		Updated:        now,
	}
	err := router.DBClient.CreateUser(r.Context(), user)
	if err != nil {
		genericResponse(w, http.StatusInternalServerError)
		return
	}

	// Generate a token for the user from the session manager
	token := router.Authorizer.GenSessionToken(user.ID)
	cookie := GetLoginCookie(token)

	http.SetCookie(w, cookie)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}

func (Router) GetSignup(w http.ResponseWriter, r *http.Request) {
	staticRender(render.Signup(), w, r)
	return
}

func (Router) GetLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   auth.SESS_COOKIE,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}

func (router *Router) GetDashboard(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	pages, err := router.DBClient.ReadPagesByUserId(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	staticRender(render.Dashboard(user, pages), w, r)
}

type GetPagesRequest struct {
	Page string `query:"p"`
}

func (router *Router) GetPages(w http.ResponseWriter, r *http.Request) {
	req := requestBind[GetPagesRequest](w, r)
	if req == nil {
		return
	}

	user := dbqueries.GetUser(r.Context())
	page, err := strconv.Atoi(req.Page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pages, err := router.DBClient.ReadPagesByUserId(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// pageinate for fun
	pages = pages[(page-1)*render.PAGE_SIZE : page*render.PAGE_SIZE]
	staticRender(render.PageList(pages, page), w, r)
}

func (router *Router) DeletePages(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())

	n, err := router.DBClient.DeletePagesByUserId(r.Context(), user.ID)
	if err != nil {
		genericResponse(w, http.StatusInternalServerError)
		return
	}
	staticRender(render.SavePageSuccess(fmt.Sprintf("Deleted %d pages", n)), w, r)

}

type GetPageRequest struct {
	PageID string `param:"page_id"`
}

func (router *Router) GetPage(w http.ResponseWriter, r *http.Request) {
	req := requestBind[GetPageRequest](w, r)
	if req == nil {
		return
	}

	user := dbqueries.GetUser(r.Context())
	page, err := router.DBClient.ReadPageById(r.Context(), req.PageID)
	if err != nil {
		http.Error(w, "Failed to get page id", http.StatusInternalServerError)
		return
	}
	if page.UserID != user.ID {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	_, err = router.DBClient.DeletePageById(r.Context(), req.PageID)
	if err != nil {
		genericResponse(w, http.StatusInternalServerError)
		return
	}
	staticRender(render.PageCard(&page), w, r)
}

type PostPageRequest struct {
	Url string `form:"url"`
}

func (router *Router) PostPage(w http.ResponseWriter, r *http.Request) {
	req := requestBind[PostPageRequest](w, r)
	if req == nil {
		return
	}

	user := dbqueries.GetUser(r.Context())
	url := req.Url
	if url == "" {
		staticRender(render.SavePageError("URL field must be present"), w, r)
		return
	}

	now := time.Now()
	page := dbqueries.Page{
		ID:      tools.GenerateNewId(20),
		UserID:  user.ID,
		Url:     url,
		Created: now,
		Updated: now,
	}
	err := router.DBClient.CreatePage(r.Context(), dbqueries.CreatePageParams{
		ID:      page.ID,
		UserID:  page.UserID,
		Url:     page.Url,
		Created: page.Created,
		Updated: page.Updated,
	})
	if err != nil {
		staticRender(render.SavePageError(err.Error()), w, r)
	}

	go func(cli *dbqueries.Queries) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := preview.FetchPreview(ctx, &page)
		if err != nil {
			return
		}
		err = cli.UpsertPage(ctx, dbqueries.UpsertPageParams{
			ID:                  page.ID,
			UserID:              page.UserID,
			Url:                 page.Url,
			Title:               page.Title,
			Description:         page.Description,
			ImageUrl:            page.ImageUrl,
			ReadabilityStatus:   page.ReadabilityStatus,
			ReadabilityTaskData: page.ReadabilityTaskData,
			IsReadable:          page.IsReadable,
			Created:             page.Created,
			Updated:             page.Updated,
		})
		if err != nil {
			return
		}
	}(router.DBClient)

	staticRender(render.PageCard(&page), w, r)
	return
}

func (router *Router) DeletePage(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	page, err := router.DBClient.ReadPageById(r.Context(), r.PathValue("page_id"))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if !router.Authorizer.ValUserAgainstPage(user.ID, page.UserID) {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	router.DBClient.DeletePageById(r.Context(), page.ID)
	w.WriteHeader(http.StatusOK)
}

func (router *Router) GetAccountPage(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	staticRender(render.AccountPage(user), w, r)
	return
}

func (router *Router) PutAccount(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	form := new(AccountData)
	err := router.DBClient.UpdateUserSubscription(r.Context(), dbqueries.UpdateUserSubscriptionParams{
		Subscribed: form.Subscribed == "on",
		ID:         user.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (router *Router) GetShortcutToken(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	token := router.Authorizer.GenShortcutToken(user.ID)
	err := router.DBClient.UpdateUserShortcutToken(r.Context(), dbqueries.UpdateUserShortcutTokenParams{
		ShortcutToken: token,
		ID:            user.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", token)
}
