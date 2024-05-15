package router

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/mr55p-dev/pagemail/internal/render"
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

func isHtmx(r *http.Request) bool {
	return r.Header.Get("Hx-Request") == "true"
}

func errorResponse(w http.ResponseWriter, r *http.Request, detail string, status int) {
	if isHtmx(r) {
		w.WriteHeader(status)
		componentRender(render.ErrorBox("Error", detail), w, r)
	} else {
		http.Error(w, fmt.Sprintf("Error: %s", detail), status)
	}
}

func textRender(message string, w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, message)
	if err != nil {
		logger.WithRequest(r).WithError(err).ErrorCtx(r.Context(), "Failed to write to request")
		http.Error(w, "Error rendering response", http.StatusInternalServerError)
	}
}

func componentRender(component templ.Component, w http.ResponseWriter, r *http.Request) {
	err := component.Render(r.Context(), w)
	if err != nil {
		logger.WithRequest(r).WithError(err).ErrorCtx(r.Context(), "Error rendering response")
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
