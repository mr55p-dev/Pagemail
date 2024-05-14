package router

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gorilla/sessions"
	"github.com/mr55p-dev/pagemail/internal/assets"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
)

var logger = logging.NewLogger("router")

const NumKeyBytes = 16

type Router struct {
	DBClient *dbqueries.Queries
	Sessions sessions.Store
	Conn     *sql.DB
	Mux      http.Handler
}

func getUserMux(router *Router) http.Handler {
	userMux := http.NewServeMux()
	userMux.HandleFunc("GET /logout", router.GetLogout)
	userMux.HandleFunc("GET /account", router.GetAccountPage)
	userMux.HandleFunc("PUT /account", router.PutAccount)
	userMux.HandleFunc("GET /token/shortcut", router.GetShortcutToken)
	return middlewares.WithMiddleware(
		http.StripPrefix("/user", userMux),
		middlewares.ProtectRoute,
	)

}

func getPagesMux(router *Router) http.Handler {
	pagesMux := http.NewServeMux()
	pagesMux.HandleFunc("GET /{page_id}", router.GetPage)
	pagesMux.HandleFunc("GET /dashboard", router.GetDashboard)
	pagesMux.HandleFunc("POST /", router.PostPage)
	pagesMux.HandleFunc("DELETE /", router.DeletePages)
	pagesMux.HandleFunc("DELETE /{page_id}", router.DeletePage)
	return middlewares.WithMiddleware(
		http.StripPrefix("/pages", pagesMux),
		middlewares.ProtectRoute,
	)
}

func getAssestMux(env Env) http.Handler {
	var fileHandler http.Handler
	switch env {
	case ENV_STG, ENV_PRD:
		subdir, err := fs.Sub(assets.FS, "public")
		if err != nil {
			panic(err)
		}
		fileHandler = http.FileServerFS(subdir)
	default:
		fileHandler = http.FileServer(http.Dir("internal/assets/public/"))
	}
	return fileHandler
}

func loadCookieKey(router *Router, cfg *AppConfig) error {
	var input io.Reader
	if cfg.CookieKeyFile == "-" {
		input = rand.Reader
	} else {
		cookieDataFile, err := os.Open(cfg.CookieKeyFile)
		if err != nil {
			return fmt.Errorf("Failed to open cookie key file: %w", err)
		}
		defer cookieDataFile.Close()
		input = cookieDataFile
	}

	key := make([]byte, NumKeyBytes)
	n, err := input.Read(key)
	if err != nil {
		return err
	}
	if n < NumKeyBytes {
		return errors.New("Cookie key source has insufficient bytes")
	}
	router.Sessions = sessions.NewCookieStore(key)
	return nil
}

func loadQueries(ctx context.Context, router *Router, cfg *AppConfig) error {
	logger.DebugCtx(ctx, "Setting up db client")
	router.Conn = dbqueries.MustGetDB(ctx, cfg.DBPath)
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
