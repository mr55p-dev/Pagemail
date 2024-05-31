package router

import (
	"context"
	"database/sql"
	"io"
	"io/fs"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
	"github.com/mr55p-dev/pagemail/internal/readability"
)

var logger = logging.NewLogger("router")

const (
	URI_DASH     = "/pages/dashboard"
	URI_LINK     = "/login/link"
	SESS_G_TOKEN = "google-token"
)

type Router struct {
	db        *sql.DB
	Previewer Previewer
	Sender    mail.Sender
	Sessions  sessions.Store
	Reader    *readability.Client
	Mux       http.Handler

	host           string
	proto          string
	googleClientId string
}

type Previewer interface {
	Queue(string)
}

func New(
	ctx context.Context,
	conn *sql.DB,
	assets fs.FS,
	mailClient mail.Sender,
	previewClient Previewer,
	cookieKey io.Reader,
	googleClientId string,
	externalHost string,
	externalProto string,
	readabilityClient *readability.Client,
) (*Router, error) {
	router := &Router{}
	router.db = conn
	router.Previewer = previewClient
	router.Sender = mailClient
	router.Reader = readabilityClient
	router.googleClientId = googleClientId
	router.host = externalHost
	router.proto = externalProto

	// Load the cookie encryption key
	router.Sessions = sessions.NewCookieStore(mustReadKey(cookieKey))

	// Serve root
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/", router.GetRoot)
	rootMux.Handle("/signup", HandleMethods(map[string]http.Handler{
		http.MethodGet:  http.HandlerFunc(router.GetSignup),
		http.MethodPost: http.HandlerFunc(router.PostSignup),
	}))
	rootMux.Handle("/shortcut/page", HandleMethod(http.MethodPost,
		middlewares.WithMiddleware(
			http.HandlerFunc(router.PostPage),
			middlewares.GetShortcutLoader(router.Sessions, router.db),
		),
	))
	rootMux.Handle("/login/", getLoginMux(router))
	rootMux.Handle("/user/", getUserMux(router))
	rootMux.Handle("/pages/", getPagesMux(router))
	rootMux.Handle("/articles/", getArticlesMux(router))
	rootMux.Handle("/password-reset/", getPasswordResetMux(router))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServerFS(assets)))
	mux.Handle("/", middlewares.WithMiddleware(rootMux,
		middlewares.Recover,
		middlewares.Tracer,
		middlewares.RequestLogger,
		middlewares.GetUserLoader(router.Sessions, router.db),
	))
	router.Mux = mux
	return router, nil
}

func mustReadKey(reader io.Reader) []byte {
	if reader == nil {
		return []byte{}
	}
	key, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return key
}
