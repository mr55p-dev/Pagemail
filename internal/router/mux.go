package router

import (
	"io/fs"
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/assets"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
)

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
