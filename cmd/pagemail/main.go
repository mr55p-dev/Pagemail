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

func HandleMethod(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			handler.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
}

func HandleMethods(methods map[string]http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := methods[r.Method]; ok {
			handler.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
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

	// Serve root
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/", s.GetRoot)
	rootMux.Handle("/login", HandleMethods(map[string]http.Handler{
		http.MethodPost: http.HandlerFunc(s.PostLogin),
		http.MethodGet:  http.HandlerFunc(s.GetLogin),
	}))

	// Serve pages
	pagesMux := http.NewServeMux()
	pagesMux.HandleFunc("GET /{page_id}", s.GetPage)
	pagesMux.HandleFunc("GET /dashboard", s.GetDashboard)
	pagesMux.HandleFunc("POST /", s.PostPage)
	pagesMux.HandleFunc("POST /shortcut", s.PostPage)
	pagesMux.HandleFunc("DELETE /", s.DeletePages)
	pagesMux.HandleFunc("DELETE /{page_id}", s.DeletePage)
	rootMux.Handle("/pages/", middlewares.WithMiddleware(
		http.StripPrefix("/pages", pagesMux),
		protector.ProtectRoute(),
	))

	// Serve users
	userMux := http.NewServeMux()
	userMux.HandleFunc("GET /logout", s.GetLogout)
	userMux.HandleFunc("GET /account", s.GetAccountPage)
	userMux.HandleFunc("PUT /account", s.PutAccount)
	userMux.HandleFunc("GET /token/shortcut", s.GetShortcutToken)
	rootMux.Handle("/user/", middlewares.WithMiddleware(
		http.StripPrefix("/user", userMux),
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
		fileHandler = http.FileServerFS(subdir)
	default:
		fileHandler = http.FileServer(http.Dir("internal/assets/public/"))
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
	mux.Handle("/assets/", http.StripPrefix("/assets", fileHandler))
	mux.Handle("/", middlewares.WithMiddleware(rootMux,
		middlewares.Recover,
		middlewares.Tracer,
		middlewares.RequestLogger,
		protector.LoadUser,
	))

	logger.Info("Starting http server", "config", cfg)
	if err := http.ListenAndServe(cfg.Host, mux); err != nil {
		panic(err)
	}
}
