package router

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/mr55p-dev/pagemail/internal/request"
)

func HandleMethod(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			handler.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
}

func HandleMethods(methods map[string]http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := methods[r.Method]; ok {
			handler.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
}

func staticRender(component templ.Component, w http.ResponseWriter, r *http.Request) {
	err := component.Render(r.Context(), w)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Error rendering response")
		http.Error(w, "Error rendering response", http.StatusInternalServerError)
	}
}
func requestBind[T any](w http.ResponseWriter, r *http.Request) *T {
	out := new(T)
	err := request.Bind(out, r)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to bind request")
		http.Error(w, "Failed to bind request", http.StatusBadRequest)
		return nil
	}
	return out
}
func genericResponse(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
