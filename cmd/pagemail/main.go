package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mr55p-dev/gonk"
	"github.com/mr55p-dev/pagemail/internal/assets"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
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
	DBClient   *dbqueries.Queries
	Authorizer *auth.Authorizer
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

func getConfig() *AppConfig {
	cfg := new(AppConfig)
	err := gonk.LoadConfig(
		cfg,
		gonk.FileLoader("pagemail.yaml", true),
		gonk.EnvironmentLoader("pm"),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}

func getAwsConfig(ctx context.Context) aws.Config {
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	return awsCfg
}

func getLogger(cfg *AppConfig) *logging.Logger {
	var lvl = slog.LevelInfo
	switch strings.ToUpper(cfg.LogLevel) {
	case "DEBUG":
		lvl = slog.LevelDebug
	case "ERROR":
		lvl = slog.LevelError
	case "WARN":
		lvl = slog.LevelWarn
	case "INFO":
		lvl = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})
	logger := logging.NewLogger("main")
	logging.SetHandler(handler)
	return logger
}

func main() {
	// Load config
	ctx := context.Background()
	cfg := getConfig()
	awsCfg := getAwsConfig(ctx)
	logger := getLogger(cfg)

	// Start the clients
	logger.DebugCtx(ctx, "Setting up db client")
	dbClient, dbClose := mustGetDb(ctx, cfg.DBPath)
	defer func() {
		_ = dbClose()
	}()

	logger.DebugCtx(ctx, "Setting up auth client")
	authClient := auth.NewAuthorizer(ctx)

	user, err := dbClient.ReadUserByEmail(ctx, "")
	user.ID

	// Handle mail
	if Env(cfg.Environment) == ENV_PRD {
		logger.InfoCtx(ctx, "Starting mail job")
		mailClient := mail.NewSesMailClient(ctx, awsCfg)
		go mail.MailGo(ctx, dbClient, mailClient)
	}

	s := &Router{
		DBClient:   dbClient,
		Authorizer: authClient,
	}

	// Serve root
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/", s.GetRoot)
	rootMux.Handle("/login", HandleMethods(map[string]http.Handler{
		http.MethodGet:  http.HandlerFunc(s.GetLogin),
		http.MethodPost: http.HandlerFunc(s.PostLogin),
	}))
	rootMux.Handle("/signup", HandleMethods(map[string]http.Handler{
		http.MethodGet:  http.HandlerFunc(s.GetSignup),
		http.MethodPost: http.HandlerFunc(s.PostSignup),
	}))
	rootMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	// Serve pages
	pagesMux := http.NewServeMux()
	pagesMux.HandleFunc("GET /{page_id}", s.GetPage)
	pagesMux.HandleFunc("GET /dashboard", s.GetDashboard)
	pagesMux.HandleFunc("POST /", s.PostPage)
	pagesMux.HandleFunc("DELETE /", s.DeletePages)
	pagesMux.HandleFunc("DELETE /{page_id}", s.DeletePage)
	rootMux.Handle("/pages/", middlewares.WithMiddleware(
		http.StripPrefix("/pages", pagesMux),
		middlewares.ProtectRoute,
	))

	rootMux.Handle("/shortcut/page", HandleMethod(http.MethodPost,
		middlewares.WithMiddleware(
			http.HandlerFunc(s.PostPage),
			middlewares.GetShortcutLoader(authClient, dbClient),
		),
	))

	// Serve users
	userMux := http.NewServeMux()
	userMux.HandleFunc("GET /logout", s.GetLogout)
	userMux.HandleFunc("GET /account", s.GetAccountPage)
	userMux.HandleFunc("PUT /account", s.PutAccount)
	userMux.HandleFunc("GET /token/shortcut", s.GetShortcutToken)
	rootMux.Handle("/user/", middlewares.WithMiddleware(
		http.StripPrefix("/user", userMux),
		middlewares.ProtectRoute,
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

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", fileHandler))
	mux.Handle("/", middlewares.WithMiddleware(rootMux,
		middlewares.Recover,
		middlewares.Tracer,
		middlewares.RequestLogger,
		middlewares.GetUserLoader(authClient, dbClient),
	))

	logger.Info("Starting http server", "config", cfg)
	if err := http.ListenAndServe(cfg.Host, mux); err != nil {
		panic(err)
	}
}
