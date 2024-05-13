package router

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/mr55p-dev/pagemail/internal/assets"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
)

var logger = logging.NewLogger("router")

type Router struct {
	DBClient   *dbqueries.Queries
	Authorizer *auth.Authorizer
	Mux        http.Handler
}

func New(ctx context.Context, cfg *AppConfig, awsCfg aws.Config) (*Router, error) {
	// Start the clients
	logger.DebugCtx(ctx, "Setting up db client")
	dbClient, dbClose := dbqueries.MustGetQueries(ctx, cfg.DBPath)
	defer func() {
		_ = dbClose()
	}()

	logger.DebugCtx(ctx, "Setting up auth client")
	authClient := auth.NewAuthorizer(ctx)

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
	s.Mux = mux
	return s, nil
}
