package readings

import (
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/middlewares"
)

func NewRoutes() http.Handler {
	mux := http.NewServeMux()
	return middlewares.WithMiddleware(
		http.StripPrefix("/articles/", mux),
		middlewares.ProtectRoute,
	)
}
