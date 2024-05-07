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
	"github.com/mr55p-dev/pagemail/internal/assets"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
	"github.com/mr55p-dev/pagemail/internal/timer"
)

var logger = logging.NewLogger("routes")

type Env string
type ContentType string

const (
	ENV_DEV Env = "dev"
	ENV_STG Env = "stg"
	ENV_PRD Env = "prd"

	CONTENT_ANY   ContentType = "*/*"
	CONTENT_HTML  ContentType = "text/html"
	CONTENT_JSON  ContentType = "text/json"
	CONTENT_PLAIN ContentType = "text/plain"
)

type Router struct {
	DBClient   *db.Client
	Authorizer *auth.Authorizer
	MailClient mail.Sender
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
		Level: ParseLogLvl(cfg.LogLevel),
	})
	logger := logging.NewLogger("main")
	logging.SetHandler(handler)

	logger.Info("Setting up db client")
	// Start the clients
	ctx := context.Background()
	dbClient := db.NewClient(cfg.DBPath, nil)
	defer dbClient.Close()

	logger.Info("Setting up auth client")
	authClient := auth.NewAuthorizer(ctx)

	logger.Info("Setting up mail client")
	var mailClient mail.Sender
	switch Env(cfg.Environment) {
	case ENV_PRD, ENV_STG:
		logger.Debug("Using SES mail client")
		mailClient = mail.NewSesMailClient(ctx)
	default:
		logger.Debug("Using test mail client")
		mailClient = &mail.TestClient{}
	}

	s := &Router{
		DBClient:   dbClient,
		Authorizer: authClient,
		MailClient: mailClient,
	}

	protector := middlewares.NewProtector(
		authClient,
		dbClient,
		logging.NewLogger("protector"),
	)

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("GET /", s.GetRoot)

	httpMux.HandleFunc("GET /login", s.GetLogin)
	httpMux.HandleFunc("POST /login", s.PostLogin)
	httpMux.HandleFunc("GET /signup", s.GetSignup)
	httpMux.HandleFunc("POST /signup", s.PostSignup)

	// Serve pages
	pagesMux := http.NewServeMux()
	pagesMux.HandleFunc("GET /", s.GetPages)
	pagesMux.HandleFunc("GET /:page_id", s.GetPage)
	pagesMux.HandleFunc("GET /dashboard", s.GetDashboard)
	pagesMux.HandleFunc("POST /", s.PostPage)
	pagesMux.HandleFunc("POST /shortcut", s.PostPage)
	pagesMux.HandleFunc("DELETE /", s.DeletePages)
	pagesMux.HandleFunc("DELETE /:page_id", s.DeletePage)
	pagesMux.Handle("GET POST DELETE /pages", WithMiddleware(
		http.StripPrefix("/pages", pagesMux),
		protector.ProtectRoute(),
	))

	// Serve users
	userMux := http.NewServeMux()
	userMux.HandleFunc("GET /logout", hut.NewHandler(s.GetLogout, protector.ProtectRoute()))
	userMux.HandleFunc("GET /account", hut.NewHandler(s.GetAccountPage, protector.ProtectRoute()))
	userMux.HandleFunc("PUT /account", hut.NewHandler(s.PutAccount, protector.ProtectRoute()))
	userMux.HandleFunc("GET /token/shortcut", hut.NewHandler(s.GetShortcutToken, protector.ProtectRoute()))
	userMux.Handle("GET POST DELETE /users", WithMiddleware(
		http.StripPrefix("/users", userMux),
		protector.ProtectRoute(),
	))

	// Serve static assets
	var fileHandler http.Handler
	switch Env(cfg.Environment) {
	case ENV_STG, ENV_PRD:
		subdir, err := fs.Sub(assets.FS, "public")
		if err != nil {
			panic(err)
		}
		fileHandler = http.StripPrefix("/assets", http.FileServerFS(subdir))
	default:
		fileHandler = http.StripPrefix("/assets", http.FileServer(http.Dir("internal/assets/public/")))

	}

	// Start the background timer
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, time.Local)
	timer := timer.NewCronTimer(time.Hour*24, start)
	defer timer.Stop()
	go func() {
		for now := range timer.T {
			slog.Info("Starting mail digest", "time", now.Format(time.Stamp))
			ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
			err := mail.DoDigestJob(ctx, dbClient, mailClient)
			cancel()
			if err != nil {
				slog.ErrorContext(ctx, "Failed to send digest", "error", err.Error())
			}
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/assets/", fileHandler)
	mux.Handle("/", WithMiddleware(httpMux,
		Recover,
		Tracer,
		RequestLogger,
		protector.LoadUser,
	))

	logger.Info("Starting http server", "config", cfg)
	if err := http.ListenAndServe(cfg.Host, mux); err != nil {
		panic(err)
	}
}
