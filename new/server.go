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

type DataPages struct{}

func (Router) GetRoot(c echo.Context) error {
	return render.RenderTempate("index", c.Response(), &DataIndex{false})
}

func (Router) GetLogin(c echo.Context) error {
	return render.RenderTempate("login", c.Response(), nil)
}

func (s *Router) PostLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	if !s.Authorizer.ValidateUser(email, password) {
		return c.String(http.StatusBadRequest, "Invalid username or password")
	}

	c.Response().Header().Set("Location", "/pages")
	c.SetCookie(&http.Cookie{
		Name:  "Authorization",
		Value: "1234",
		Path:  "/",
	})
	return c.NoContent(http.StatusSeeOther)
}

func (Router) GetLogout(c echo.Context) error {
	c.Response().Header().Set("Location", "/login")
	c.SetCookie(&http.Cookie{
		Name:   "Authorization",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return c.NoContent(http.StatusSeeOther)
}

func (Router) GetPages(c echo.Context) error {
	return render.RenderTempate("pages", c.Response(), &DataPages{})
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

	e := echo.New()
	dbClient := db.NewClient(dbLogger)
	defer dbClient.Close()

	authClient := auth.NewAuthorizer(dbClient, authLogger)

	s := &Router{
		DBClient:   dbClient,
		Authorizer: authClient,
	}

	e.Use(middlewares.GetLoggingMiddleware(reqLogger))

	e.GET("/", s.GetRoot)

	e.GET("/login", s.GetLogin)
	e.GET("/logout", s.GetLogout)
	e.POST("/login", s.PostLogin)

	e.GET("/pages", s.GetPages)

	if err := e.Start(":8080"); err != nil {
		rootLogger.Error().Msg(err.Error())
	}
}
