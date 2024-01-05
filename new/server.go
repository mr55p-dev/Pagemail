package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/auth"
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/middlewares"
	"github.com/mr55p-dev/pagemail/pkg/render"
	"github.com/rs/zerolog"
)

type Router struct {
	DBClient   db.AbsClient
	Authorizer auth.AbsAuthorizer
}

type DataIndex struct {
	IsUser bool
}

type DataPages struct {
	Pages []db.Page
}

func (Router) GetRoot(c echo.Context) error {
	return render.RenderTempate("index", c.Response(), &DataIndex{false})
}

func (Router) GetLogin(c echo.Context) error {
	return render.RenderTempate("login", c.Response(), nil)
}

func (s *Router) PostLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	user, err := s.DBClient.ReadUserByEmail(email)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if !s.Authorizer.ValidateUser(email, password, user) {
		return c.String(http.StatusBadRequest, "Invalid username or password")
	}

	sess := s.Authorizer.GetToken(user.Id)

	c.Response().Header().Set("Location", "/pages")
	c.SetCookie(&http.Cookie{
		Name:  auth.SESS_COOKIE,
		Value: sess,
		Path:  "/",
	})
	return c.NoContent(http.StatusSeeOther)
}

func (Router) GetSignup(c echo.Context) error {
	return render.RenderTempate("signup", c.Response(), nil)
}

func (s *Router) PostSignup(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	passwordRepeat := c.FormValue("password-repeat")
	if password != passwordRepeat {
		return c.String(http.StatusBadRequest, "Passwords do not match")
	}

	token, err := s.Authorizer.SignupNewUser(email, password, username)
	if err != nil {
		return c.String(http.StatusBadRequest, "Something went wrong")
	}

	c.Response().Header().Set("Location", "/pages")
	c.SetCookie(&http.Cookie{
		Name:  auth.SESS_COOKIE,
		Value: token,
		Path:  "/",
	})
	return c.NoContent(http.StatusSeeOther)
}

func (Router) GetLogout(c echo.Context) error {
	c.Response().Header().Set("Location", "/login")
	c.SetCookie(&http.Cookie{
		Name:   auth.SESS_COOKIE,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return c.NoContent(http.StatusSeeOther)
}

func (r *Router) GetPages(c echo.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	pages, err := r.DBClient.ReadPagesByUserId(user.Id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return render.RenderTempate("pages", c.Response(), &DataPages{pages})
}

func (r *Router) PostPage(c echo.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok || user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	url := c.FormValue("url")
	page := db.NewPage(user.Id, url)

	if err := r.DBClient.CreatePage(page); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (Router) TestTmpl(c echo.Context) error {
	cmp := render.Hello("ellis")
	return cmp.Render(c.Request().Context(), c.Response())
}

func main() {
	// Logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logOut := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s |", i)
		},
		FormatCaller: func(i interface{}) string {
			return filepath.Base(fmt.Sprintf("%s", i))
		},
		PartsExclude: []string{
			zerolog.TimestampFieldName,
		},
	}

	rootLogger := zerolog.New(logOut).With().Timestamp().Caller().Logger()
	reqLogger := rootLogger.With().Str("service", "root").Logger()
	authLogger := rootLogger.With().Str("service", "auth").Logger()
	dbLogger := rootLogger.With().Str("service", "db").Logger()
	protectionLogger := rootLogger.With().Str("service", "protection").Logger()

	e := echo.New()
	dbClient := db.NewClient(dbLogger)
	defer dbClient.Close()

	authClient := auth.NewAuthorizer(dbClient, authLogger)

	s := &Router{
		DBClient:   dbClient,
		Authorizer: authClient,
	}

	protected := middlewares.GetProtectedMiddleware(protectionLogger, authClient, dbClient)
	e.Use(middlewares.GetLoggingMiddleware(reqLogger))

	e.GET("/", s.GetRoot)

	e.GET("/login", s.GetLogin)
	e.GET("/logout", s.GetLogout, protected)
	e.POST("/login", s.PostLogin)

	e.GET("/signup", s.GetSignup)
	e.POST("/signup", s.PostSignup)

	e.GET("/pages", s.GetPages, protected)
	e.POST("/page", s.PostPage, protected)

	if err := e.Start(":8080"); err != nil {
		rootLogger.Error().Msg(err.Error())
	}
}
