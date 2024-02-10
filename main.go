package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mr55p-dev/pagemail/docs"
	"github.com/mr55p-dev/pagemail/pkg/auth"
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/logging"
	"github.com/mr55p-dev/pagemail/pkg/mail"
	"github.com/mr55p-dev/pagemail/pkg/middlewares"
	"github.com/mr55p-dev/pagemail/pkg/preview"
	"github.com/mr55p-dev/pagemail/pkg/render"
	"github.com/robfig/cron/v3"
	"github.com/swaggo/echo-swagger"
)

type Env string
type Mode string
type ContentType string

const (
	ENV_DEV Env = "dev"
	ENV_STG Env = "stg"
	ENV_PRD Env = "prd"

	MODE_LOCAL   Mode = "local"
	MODE_RELEASE Mode = "release"

	CONTENT_ANY   ContentType = "*/*"
	CONTENT_HTML  ContentType = "text/html"
	CONTENT_JSON  ContentType = "text/json"
	CONTENT_PLAIN ContentType = "text/plain"
)

//go:embed public
var public embed.FS
var log logging.Log

type Router struct {
	DBClient   *db.Client
	Authorizer auth.Authorizer
	MailClient mail.MailClient
}

func SetRedirect(c echo.Context, dest string) {
	c.Response().Header().Set("HX-Location", dest)
}

func SetLoginCookie(c echo.Context, val string) {
	c.SetCookie(&http.Cookie{
		Name:     auth.SESS_COOKIE,
		Value:    val,
		Path:     "/",
		MaxAge:   864000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func LoadUser(c echo.Context) *db.User {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		panic(fmt.Errorf("Could not load request user"))
	}
	return user
}

func MakeErrorResponse(c echo.Context, status int, err error) error {
	return render.ReturnRender(c, render.ErrorBox(err.Error()))
}

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
	log.ReqInfo(c, "requested login", logging.UserMail, email)
	user, err := r.DBClient.ReadUserByEmail(c.Request().Context(), email)
	if err != nil {
		log.ReqErr(c, "DB error when logging in", err)
		return MakeErrorResponse(
			c, http.StatusInternalServerError, err,
		)
	}

	if !r.Authorizer.ValCredentialsAgainstUser(email, password, user) {
		log.ReqErr(c, "Unauthorized login attempt", err, logging.UserMail, email)
		return MakeErrorResponse(c, http.StatusUnauthorized, err)
	}

	sess := r.Authorizer.GenSessionToken(user)
	log.ReqDebug(c, "Login succesful", logging.User, user)

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

	return render.ReturnRender(c, render.PageCard(page))
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
	log.ReqDebug(c, "Account page", logging.User, user)
	return render.ReturnRender(c, render.AccountPage(user))
}

func (r *Router) GetShortcutToken(c echo.Context) error {
	user := LoadUser(c)
	log.DebugContext(c.Request().Context(), "Requested new shortcut token")
	token := r.Authorizer.GenShortcutToken(user)
	user.ShortcutToken = token
	err := r.DBClient.UpdateUser(c.Request().Context(), user)
	if err != nil {
		msg := "Failed to load entries from database"
		log.ReqErr(c, msg, err, logging.User, user)
		return MakeErrorResponse(c, http.StatusInternalServerError, err)
	}
	return c.String(http.StatusOK, token)
}

func main() {
	log = logging.GetLogger("root")
	cfg := logging.Config
	log.Info("Configuring server", "config", cfg)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	ctx := context.Background()
	dbClient := db.NewClient(cfg.DBPath)
	defer dbClient.Close()

	var mailClient mail.MailClient
	var authClient auth.Authorizer
	switch Env(cfg.Env) {
	case ENV_PRD, ENV_STG:
		log.Info("Using real auth")
		authClient = auth.NewSecureAuthorizer(ctx, dbClient)
		mailClient = mail.NewSesMailClient(ctx)
	default:
		log.Info("Using test auth")
		authClient = auth.NewTestAuthorizer(cfg.TestUser)
		mailClient = &mail.TestClient{}
	}

	s := &Router{
		DBClient:   dbClient,
		Authorizer: authClient,
		MailClient: mailClient,
	}

	e.Use(middlewares.TraceMiddleware)
	e.Use(middlewares.GetLoggingMiddleware)
	if authClient == nil {
		panic("nil auth client")
	}
	tryLoadUser := middlewares.GetProtectedMiddleware(authClient, dbClient, false)
	protected := middlewares.GetProtectedMiddleware(authClient, dbClient, true)
	shortcut := middlewares.GetShortcutProtected(authClient, dbClient)

	switch Mode(cfg.Mode) {
	case MODE_RELEASE:
		fs := echo.MustSubFS(public, "public")
		e.StaticFS("/assets", fs)
	default:
		e.Static("/assets", "public")
	}

	e.GET("/", s.GetRoot, tryLoadUser)

	e.GET("/login", s.GetLogin)
	e.POST("/login", s.PostLogin)

	e.GET("/signup", s.GetSignup)
	e.POST("/signup", s.PostSignup)

	e.GET("/logout", s.GetLogout, protected)

	e.GET("/account", s.GetAccountPage, protected)
	e.GET("/shortcut-token", s.GetShortcutToken, protected)

	e.GET("/dashboard", s.GetDashboard, protected)
	e.GET("/pages", s.GetPages, protected)
	e.DELETE("/pages", s.DeletePages, protected)

	e.GET("/page/:page_id", s.GetPage, protected)
	e.DELETE("/page/:page_id", s.DeletePage, protected)
	e.POST("/page", s.PostPage, protected)

	e.POST("/shortcut/page", s.PostPage, shortcut)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// e.GET("/:id/pages/listen", s.ListenPages, protected)
	// e.GET("/test", s.TestUpdate, protected)

	cr := cron.New()
	cr.AddFunc(
		"0 7 * * *",
		func() {
			ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
			defer cancel()
			mail.DoDigestJob(ctx, dbClient, mailClient)
		},
	)

	if err := e.Start(fmt.Sprintf("127.0.0.1:%s", cfg.Port)); err != nil {
		log.Err("Server exited with error", err)
	}
}
