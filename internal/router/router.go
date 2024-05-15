package router

import (
	"context"
	"database/sql"
	"io"
	"io/fs"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
	"github.com/mr55p-dev/pagemail/internal/preview"
)

var logger = logging.NewLogger("router")

type Router struct {
	DBClient  *dbqueries.Queries
	Previewer *preview.Client
	Sessions  sessions.Store
	Mux       http.Handler
}

func New(
	ctx context.Context,
	conn *dbqueries.Queries,
	assets fs.FS,
	mailClient mail.Sender,
	previewClient *preview.Client,
	cookieKey io.Reader,
) (*Router, error) {
	router := &Router{}
	router.DBClient = conn
	router.Previewer = previewClient

	// Load the cookie encryption key
	err := loadCookieKey(router, cookieKey)
	if err != nil {
		return nil, err
	}

	// Load the mail client
	go mail.MailGo(ctx, router.DBClient, mailClient)

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

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServerFS(assets)))
	mux.Handle("/", middlewares.WithMiddleware(rootMux,
		middlewares.Recover,
		middlewares.Tracer,
		middlewares.RequestLogger,
		middlewares.GetUserLoader(router.Sessions, router.DBClient),
	))
	router.Mux = mux
	return router, nil
}
