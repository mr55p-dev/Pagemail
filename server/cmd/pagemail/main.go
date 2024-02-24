package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	configLoader "github.com/mr55p-dev/config-loader"
	"github.com/mr55p-dev/pagemail/internal/assets"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
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

type Router struct {
	DBClient   *db.Client
	Authorizer auth.Authorizer
	MailClient mail.MailClient
}

type AccountForm struct {
	Subscribed string `form:"email-list"`
}

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	baseLog := logging.Logger{
		Logger: slog.New(handler),
	}
	log := baseLog.With("module", "main")

	cfg := new(AppConfig)
	err := configLoader.LoadConfig(
		cfg,
		configLoader.EnvironmentLoader("pm"),
		configLoader.FileLoader("config.yaml", true),
	)
	if err != nil {
		log.Error("failed to load config", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	ctx := context.Background()
	dbClient := db.NewClient(cfg.DBPath, baseLog.Module("db"))
	defer dbClient.Close()

	var mailClient mail.MailClient
	var authClient auth.Authorizer
	switch Env(cfg.Environment) {
	case ENV_PRD, ENV_STG:
		tokens, err := dbClient.ReadUserShortcutTokens(ctx)
		if err != nil {
			panic(err.Error())
		}
		authClient = auth.NewSecureAuthorizer(ctx, tokens...)
		mailClient = mail.NewSesMailClient(ctx)
	default:
		authClient = auth.NewTestAuthorizer()
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

	switch Env(cfg.Environment) {
	case ENV_STG, ENV_PRD:
		fs := echo.MustSubFS(assets.FS, "public")
		e.StaticFS("/assets", fs)
	default:
		e.Static("/assets", "server/internal/assets/public")
	}

	e.GET("/", s.GetRoot, tryLoadUser)

	e.GET("/login", s.GetLogin)
	e.POST("/login", s.PostLogin)
	e.GET("/signup", s.GetSignup)
	e.POST("/signup", s.PostSignup)
	e.GET("/logout", s.GetLogout, protected)

	e.GET("/dashboard", s.GetDashboard, protected)
	e.GET("/pages", s.GetPages, protected)
	e.DELETE("/pages", s.DeletePages, protected)
	e.GET("/page/:page_id", s.GetPage, protected)
	e.DELETE("/page/:page_id", s.DeletePage, protected)
	e.POST("/page", s.PostPage, protected)

	e.GET("/account", s.GetAccountPage, protected)
	e.PUT("/account", s.PutAccount, protected)

	e.GET("/shortcut-token", s.GetShortcutToken, protected)
	e.POST("/shortcut/page", s.PostPage, shortcut)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	cr := cron.New()
	cr.AddFunc(
		"0 7 * * *",
		func() {
			ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
			defer cancel()
			mail.DoDigestJob(ctx, dbClient, mailClient)
		},
	)

	if err := e.Start(fmt.Sprintf("127.0.0.1:%s", "8080")); err != nil {
	}
}
