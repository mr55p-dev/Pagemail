package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/auth"
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/logging"
	"github.com/mr55p-dev/pagemail/pkg/mail"
	"github.com/mr55p-dev/pagemail/pkg/middlewares"
	"github.com/mr55p-dev/pagemail/pkg/preview"
	"github.com/mr55p-dev/pagemail/pkg/render"
	"github.com/robfig/cron/v3"
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

func GetAccept(c echo.Context) ContentType {
	accept := c.Request().Header.Get("Accept")
	return ContentType(accept)
}

func (Router) GetRoot(c echo.Context) error {
	var user *db.User
	if c.Get("user") != nil {
		user = c.Get("user").(*db.User)
	}
	return render.ReturnRender(c, render.Index(user))
}

func (Router) GetLogin(c echo.Context) error {
	return render.ReturnRender(c, render.Login())
}

func (s *Router) PostLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	log.ReqInfo(c, "requested login", logging.UserMail, email)
	user, err := s.DBClient.ReadUserByEmail(c.Request().Context(), email)
	if err != nil {
		log.ReqErr(c, "DB error when logging in", err)
		return c.NoContent(http.StatusBadRequest)
	}

	if !s.Authorizer.ValCredentialsAgainstUser(email, password, user) {
		log.ReqErr(c, "Unauthorized login attempt", err, logging.UserMail, email)
		return c.String(http.StatusBadRequest, "Invalid username or password")
	}

	sess := s.Authorizer.GenSessionToken(user)
	log.ReqDebug(c, "Login succesfull", logging.User, user)

	c.Response().Header().Set("HX-Location", fmt.Sprintf("/%s/pages", user.Id))
	c.SetCookie(&http.Cookie{
		Name:  auth.SESS_COOKIE,
		Value: sess,
		Path:  "/",
	})
	return c.NoContent(http.StatusOK)
}

func (Router) GetSignup(c echo.Context) error {
	return render.ReturnRender(c, render.Signup())
}

func (s *Router) PostSignup(c echo.Context) error {
	// Read the form requests
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	passwordRepeat := c.FormValue("password-repeat")
	if password != passwordRepeat {
		return c.String(http.StatusBadRequest, "Passwords do not match")
	}

	// Generate a new user
	user := db.NewUser(email, s.Authorizer.GenPasswordHash(password))
	user.Username = username
	err := s.DBClient.CreateUser(c.Request().Context(), user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Something went wrong")
	}

	// Generate a token for the user from the session manager
	token := s.Authorizer.GenSessionToken(user)

	c.Response().Header().Set("HX-Location", fmt.Sprintf("/%s/pages", user.Id))
	c.SetCookie(&http.Cookie{
		Name:  auth.SESS_COOKIE,
		Value: token,
		Path:  "/",
	})
	return c.NoContent(http.StatusOK)
}

func (Router) GetLogout(c echo.Context) error {
	c.Response().Header().Set("Location", "/login")
	c.SetCookie(&http.Cookie{
		Name:   auth.SESS_COOKIE,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	c.Response().Header().Set("HX-Location", "/")
	return c.NoContent(http.StatusOK)
}

func (r *Router) GetPages(c echo.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	pages, err := r.DBClient.ReadPagesByUserId(c.Request().Context(), user.Id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return render.ReturnRender(c, render.PageView(user, pages))
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

	if ShouldRender(c) {
		return render.ReturnRender(c, render.SavePageSuccess(fmt.Sprintf("Deleted %d records", n)))
	} else {
		return c.String(http.StatusOK, fmt.Sprintf("Deleted %d records", n))
	}

}

func (r *Router) GetPage(c echo.Context) error {
	user := c.Get("user").(*db.User)
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
	render.PageElementComponent(page).Render(
		c.Request().Context(),
		c.Response().Writer,
	)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func ShouldRender(c echo.Context) bool {
	return c.Request().Header.Get("Accept") == "*/*"
}

func (r *Router) PostPage(c echo.Context) error {

	user, ok := c.Get("user").(*db.User)
	if !ok || user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	url := c.FormValue("url")
	if url == "" {
		if ShouldRender(c) {
			return render.ReturnRender(c, render.SavePageError("URL field must be present"))
		} else {
			return c.String(http.StatusBadRequest, "URL field must be present")
		}
	}

	page := db.NewPage(user.Id, url)
	if err := r.DBClient.CreatePage(c.Request().Context(), page); err != nil {
		if ShouldRender(c) {
			return render.ReturnRender(c, render.SavePageError(err.Error()))
		} else {
			return c.String(http.StatusInternalServerError, err.Error())
		}
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

	if ShouldRender(c) {
		return render.ReturnRender(c, render.PageElementComponent(page))
	} else {
		return c.NoContent(http.StatusCreated)
	}
}

func (r *Router) DeletePage(c echo.Context) error {
	user := c.Get("user").(*db.User)
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

type Error struct {
	Message string `json:"message"`
}

func MakeErrorResponse(c echo.Context, status int, message string, err error, component templ.Component) error {
	switch GetAccept(c) {
	case CONTENT_PLAIN:
		return c.String(status, message)
	case CONTENT_JSON:
		return c.JSON(status, Error{message})
	default:
		return render.ReturnRender(c, component)
	}
}

func (r *Router) GetAccountPage(c echo.Context) error {
	user := c.Get("user").(*db.User)
	log.ReqDebug(c, "Account page", logging.User, user)
	return render.ReturnRender(c, render.AccountPage(user))
}

func (r *Router) GetShortcutToken(c echo.Context) error {
	user := c.Get("user").(*db.User)
	log.DebugContext(c.Request().Context(), "Requested new shortcut token")
	token := r.Authorizer.GenShortcutToken(user)
	user.ShortcutToken = token
	err := r.DBClient.UpdateUser(c.Request().Context(), user)
	if err != nil {
		msg := "Failed to load entries from database"
		log.ReqErr(c, msg, err, logging.User, user)
		return MakeErrorResponse(
			c, http.StatusInternalServerError,
			msg, err, render.ErrorBox(msg),
		)
	}
	return c.String(http.StatusOK, token)
}

// func (r *Router) ListenPages(c echo.Context) error {
// 	user, ok := c.Get("user").(*db.User)
// 	if !ok || user == nil {
// 		return c.NoContent(http.StatusInternalServerError)
// 	}
//
// 	// create listener to updates
// 	events := make(chan db.Event[db.Page])
//
// 	r.DBClient.AddPageListener(user.Id, events)
//
// 	c.Response().WriteHeader(http.StatusOK)
// 	c.Response().Header().Add("Content-Type", "text/event-stream")
// 	c.Response().Flush()
//
// 	defer r.DBClient.RemovePageListener(user.Id)
// 	for event := range events {
// 		fmt.Println("Event", event)
// 		c.Response().Write([]byte("<p>Message sent</p>\n\n"))
// 		c.Response().Flush()
// 	}
//
// 	return nil
// }

// func (r *Router) TestUpdate(c echo.Context) error {
// 	user, _ := c.Get("user").(*db.User)
// 	page := db.NewPage(user.Id, "https://example.com")
// 	page.Id = "123456"
// 	err := r.DBClient.UpsertPage(c.Request().Context(), page)
// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	return c.String(200, "Updated")
// }

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
	protected := middlewares.GetProtectedMiddleware(authClient, dbClient)

	switch Mode(cfg.Mode) {
	case MODE_RELEASE:
		fs := echo.MustSubFS(public, "public")
		e.StaticFS("/assets", fs)
	default:
		e.Static("/assets", "public")
	}

	e.GET("/", s.GetRoot)

	e.GET("/login", s.GetLogin)
	e.POST("/login", s.PostLogin)

	e.GET("/signup", s.GetSignup)
	e.POST("/signup", s.PostSignup)

	e.GET("/:id/logout", s.GetLogout, protected)

	e.GET("/:id/account", s.GetAccountPage, protected)
	e.GET("/:id/shortcut-token", s.GetShortcutToken, protected)

	e.GET("/:id/pages", s.GetPages, protected)
	e.DELETE("/:id/pages", s.DeletePages, protected)

	e.GET("/:id/page/:page_id", s.GetPage, protected)
	e.DELETE("/:id/page/:page_id", s.DeletePage, protected)
	e.POST("/:id/page", s.PostPage, protected)

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

	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Err("Server exited with error", err)
	}
}
