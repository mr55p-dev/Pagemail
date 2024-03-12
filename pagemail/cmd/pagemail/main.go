package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mr55p-dev/gonk"
	"github.com/mr55p-dev/htmx-utils"
	hutMiddlewares "github.com/mr55p-dev/htmx-utils/pkg/middlewares"
	"github.com/mr55p-dev/pagemail/internal/assets"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
	"github.com/mr55p-dev/pagemail/internal/tools"
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

type Router struct {
	DBClient   *db.Client
	Authorizer *auth.Authorizer
	MailClient mail.MailClient
	log        *logging.Logger
}

type AccountData struct {
	Subscribed string `form:"email-list"`
}

func ParseLogLvl(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "ERROR":
		return slog.LevelError
	case "WARN":
		return slog.LevelWarn
	case "INFO":
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}

func main() {
	// Load config
	cfg := new(AppConfig)
	err := gonk.LoadConfig(
		cfg,
		gonk.FileLoader("pagemail.yaml", true),
		gonk.EnvironmentLoader("pm"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		panic(err)
	}

	// Setup logging
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     ParseLogLvl(cfg.LogLevel),
		AddSource: true,
	})
	baseLog := logging.Logger{
		Logger: slog.New(handler),
	}

	// Basic middleware
	hut.UseGlobal(
		hutMiddlewares.Recover,
		hutMiddlewares.RequestLogger(baseLog.With("module", "request logger")),
	)

	// Start the clients
	ctx := context.Background()
	dbClient := db.NewClient(cfg.DBPath, &baseLog)
	defer dbClient.Close()

	authClient := auth.NewAuthorizer(ctx)

	var mailClient mail.MailClient
	switch Env(cfg.Environment) {
	case ENV_PRD, ENV_STG:
		mailClient = mail.NewSesMailClient(ctx, logging.New(baseLog.With("package", "mail")))
	default:
		mailClient = &mail.TestClient{}
	}

	s := &Router{
		DBClient:   dbClient,
		Authorizer: authClient,
		MailClient: mailClient,
		log:        &baseLog,
	}

	hut.UseGlobal(
		hutMiddlewares.Trace(func() string {
			return tools.GenerateNewId(10)
		}),
	)
	if authClient == nil {
		panic("nil auth client")
	}

	protector := middlewares.NewProtector(
		authClient,
		dbClient,
		logging.New(baseLog.With("package", "protection middleware")),
	)
	hut.UseGlobal(protector.LoadUser)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", hut.NewHandler(s.GetRoot))

	mux.HandleFunc("GET /login", hut.NewHandler(s.GetLogin))
	mux.HandleFunc("POST /login", hut.NewBoundHandler(s.PostLogin))
	mux.HandleFunc("GET /signup", hut.NewHandler(s.GetSignup))
	mux.HandleFunc("POST /signup", hut.NewBoundHandler(s.PostSignup))
	mux.HandleFunc("GET /logout", hut.NewHandler(s.GetLogout, protector.ProtectRoute()))

	mux.HandleFunc("GET /dashboard", hut.NewHandler(s.GetDashboard, protector.ProtectRoute()))
	mux.HandleFunc("GET /pages", hut.NewBoundHandler(s.GetPages, protector.ProtectRoute()))
	mux.HandleFunc("DELETE /pages", hut.NewHandler(s.DeletePages, protector.ProtectRoute()))
	mux.HandleFunc("GET /page/:page_id", hut.NewBoundHandler(s.GetPage, protector.ProtectRoute()))
	mux.HandleFunc("DELETE /page/:page_id", hut.NewBoundHandler(s.DeletePage, protector.ProtectRoute()))
	mux.HandleFunc("POST /page", hut.NewBoundHandler(s.PostPage, protector.ProtectRoute()))

	mux.HandleFunc("GET /account", hut.NewHandler(s.GetAccountPage, protector.ProtectRoute()))
	mux.HandleFunc("PUT /account", hut.NewHandler(s.PutAccount, protector.ProtectRoute()))

	mux.HandleFunc("GET /shortcut-token", hut.NewHandler(s.GetShortcutToken, protector.ProtectRoute()))
	mux.HandleFunc("POST /shortcut/page", hut.NewBoundHandler(s.PostPage, protector.LoadFromShortcut()))

	switch Env(cfg.Environment) {
	case ENV_STG, ENV_PRD:
		subdir, err := fs.Sub(assets.FS, "public")
		if err != nil {
			panic(err)
		}
		mux.Handle("GET /assets/", http.FileServerFS(subdir))
	default:
		mux.Handle("GET /assets/", http.StripPrefix(
			"/assets",
			http.FileServer(http.Dir("pagemail/internal/assets/public/"))),
		)
	}

	mailLog := logging.New(baseLog.With("package", "mail"))
	cr := cron.New()
	cr.AddFunc(
		"0 7 * * *",
		func() {
			ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
			defer cancel()
			mail.DoDigestJob(ctx, mailLog, dbClient, mailClient)
		},
	)

	baseLog.Info("Starting http server", "config", cfg)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), mux); err != nil {
		panic(err)
	}
}
