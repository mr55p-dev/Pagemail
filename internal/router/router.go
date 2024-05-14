package router

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gorilla/sessions"
	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
)

var logger = logging.NewLogger("router")

type Router struct {
	DBClient *dbqueries.Queries
	Sessions sessions.Store
	Conn     *sql.DB
	Mux      http.Handler
}

func New(ctx context.Context, cfg *AppConfig) (*Router, error) {
	router := &Router{}

	// Load the cookie encryption key
	err := loadCookieKey(router, cfg)
	if err != nil {
		return nil, err
	}

	// Load the db queries
	err = loadQueries(ctx, router, cfg)
	if err != nil {
		return nil, err
	}

	// Load the mail client
	err = loadMailer(ctx, router, cfg)
	if err != nil {
		return nil, err
	}

	// Serve root
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/", router.GetRoot)
	rootMux.Handle("/login", HandleMethods(map[string]http.Handler{
		http.MethodGet:  http.HandlerFunc(router.GetLogin),
		http.MethodPost: http.HandlerFunc(router.PostLogin),
	}))
	rootMux.Handle("/signup", HandleMethods(map[string]http.Handler{
		http.MethodGet:  http.HandlerFunc(router.GetSignup),
		http.MethodPost: http.HandlerFunc(router.PostSignup),
	}))
	rootMux.Handle("/shortcut/page", HandleMethod(http.MethodPost,
		middlewares.WithMiddleware(
			http.HandlerFunc(router.PostPage),
			middlewares.GetShortcutLoader(router.Sessions, router.DBClient),
		),
	))
	rootMux.Handle("/user/", getUserMux(router))
	rootMux.Handle("/pages/", getPagesMux(router))

	fileHandler := getAssestMux(Env(cfg.Environment))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", fileHandler))
	mux.Handle("/", middlewares.WithMiddleware(rootMux,
		middlewares.Recover,
		middlewares.Tracer,
		middlewares.RequestLogger,
		middlewares.GetUserLoader(router.Sessions, router.DBClient),
	))
	router.Mux = mux
	return router, nil
}

func loadQueries(ctx context.Context, router *Router, cfg *AppConfig) error {
	logger.DebugCtx(ctx, "Setting up db client")
	router.Conn = db.MustConnect(ctx, cfg.DBPath)
	router.DBClient = dbqueries.New(router.Conn)
	go func() {
		<-ctx.Done()
		_ = router.Conn.Close()
	}()
	return nil
}

func loadMailer(ctx context.Context, router *Router, cfg *AppConfig) error {
	if Env(cfg.Environment) == ENV_PRD {
		awsCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			panic(err)
		}
		logger.InfoCtx(ctx, "Starting mail job")
		mailClient := mail.NewSesMailClient(ctx, awsCfg)
		go mail.MailGo(ctx, router.DBClient, mailClient)
	}
	return nil
}
