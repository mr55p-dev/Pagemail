package router

import (
	"context"
	"database/sql"
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
	Authorizer *auth.MemoryStore
	Conn       *sql.DB
	Mux        http.Handler
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

func New(ctx context.Context, cfg *AppConfig, awsCfg aws.Config) (*Router, error) {
	// Start the clients
	router := &Router{}
	logger.DebugCtx(ctx, "Setting up db client")
	router.Conn = dbqueries.MustGetDB(ctx, cfg.DBPath)
	router.DBClient = dbqueries.New(router.Conn)
	go func() {
		<-ctx.Done()
		_ = router.Conn.Close()
	}()

	logger.DebugCtx(ctx, "Setting up auth client")
	router.Authorizer = auth.NewMemoryStore(ctx)

	// Handle mail
	if Env(cfg.Environment) == ENV_PRD {
		logger.InfoCtx(ctx, "Starting mail job")
		mailClient := mail.NewSesMailClient(ctx, awsCfg)
		go mail.MailGo(ctx, router.DBClient, mailClient)
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
			middlewares.GetShortcutLoader(router.Authorizer, router.DBClient),
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
		middlewares.GetUserLoader(router.Authorizer, router.DBClient),
	))
	router.Mux = mux
	return router, nil
}
