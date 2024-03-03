package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	configLoader "github.com/mr55p-dev/config-loader"
	"github.com/mr55p-dev/go-httpit"
	httpItMiddlewares "github.com/mr55p-dev/go-httpit/pkg/middlewares"
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
	Authorizer auth.Authorizer
	MailClient mail.MailClient
}

type AccountData struct {
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

	httpit.UseGlobal(
		httpItMiddlewares.Recover,
		httpItMiddlewares.RequestLogger(baseLog.With("module", "request logger")),
	)

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
		mailClient = mail.NewSesMailClient(ctx, baseLog.Module("mail"))
	default:
		authClient = auth.NewTestAuthorizer()
		mailClient = &mail.TestClient{}
	}

	s := &Router{
		DBClient:   dbClient,
		Authorizer: authClient,
		MailClient: mailClient,
	}

	httpit.UseGlobal(
		httpItMiddlewares.Trace(func() string {
			return tools.GenerateNewId(10)
		}),
	)
	if authClient == nil {
		panic("nil auth client")
	}

	protector := middlewares.NewProtector(authClient, dbClient, baseLog.Module("protection middleware"))
	httpit.UseGlobal(protector.LoadUser)
	mux := http.NewServeMux()

	switch Env(cfg.Environment) {
	case ENV_STG, ENV_PRD:
		subdir, err := fs.Sub(assets.FS, "public")
		if err != nil {
			panic(err)
		}
		mux.Handle("GET /assets", http.FileServerFS(subdir))
	default:
		mux.Handle("GET /assets", http.FileServer(http.Dir("server/internal/assets/public")))
	}

	mux.HandleFunc("GET /", httpit.NewTemplHandler(s.GetRoot))

	mux.HandleFunc("GET /login", httpit.NewTemplHandler(s.GetLogin))
	mux.HandleFunc("POST /login", httpit.NewInHandler(s.PostLogin))
	mux.HandleFunc("GET /signup", httpit.NewTemplHandler(s.GetSignup))
	mux.HandleFunc("POST /signup", httpit.NewInHandler(s.PostSignup))
	mux.HandleFunc("GET /logout", httpit.NewHandler(s.GetLogout, protector.ProtectRoute()))

	mux.HandleFunc("GET /dashboard", httpit.NewTemplHandler(s.GetDashboard, protector.ProtectRoute()))
	mux.HandleFunc("GET /pages", httpit.NewTemplMappedHandler(s.GetPages, protector.ProtectRoute()))
	mux.HandleFunc("DELETE /pages", httpit.NewTemplHandler(s.DeletePages, protector.ProtectRoute()))
	mux.HandleFunc("GET /page/:page_id", httpit.NewTemplMappedHandler(s.GetPage, protector.ProtectRoute()))
	mux.HandleFunc("DELETE /page/:page_id", httpit.NewInHandler(s.DeletePage, protector.ProtectRoute()))
	mux.HandleFunc("POST /page", httpit.NewTemplMappedHandler(s.PostPage, protector.ProtectRoute()))

	mux.HandleFunc("GET /account", httpit.NewTemplHandler(s.GetAccountPage, protector.ProtectRoute()))
	mux.HandleFunc("PUT /account", httpit.NewTemplHandler(s.PutAccount, protector.ProtectRoute()))

	mux.HandleFunc("GET /shortcut-token", httpit.NewHandler(s.GetShortcutToken, protector.ProtectRoute()))
	mux.HandleFunc("POST /shortcut/page", httpit.NewTemplMappedHandler(s.PostPage, protector.LoadFromShortcut()))

	mailLog := baseLog.Module("mail")
	cr := cron.New()
	cr.AddFunc(
		"0 7 * * *",
		func() {
			ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
			defer cancel()
			mail.DoDigestJob(ctx, mailLog, dbClient, mailClient)
		},
	)

	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", "8080"), mux); err != nil {
		panic(err)
	}
}
