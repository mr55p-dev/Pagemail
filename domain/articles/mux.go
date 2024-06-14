package pages

import (
	"database/sql"
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/middlewares"
)

func Routes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", NewReadArticles())
	mux.Handle("POST /", NewCreateArticle(db))

	return middlewares.WithMiddleware(
		mux,
		middlewares.ProtectRoute,
	)
}
