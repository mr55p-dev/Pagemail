package httpit

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

var (
	HX_REQUEST  = "Hx-Request"
	HX_REDIRECT = "Hx-Redirect"

	globalMiddlewares []MiddlewareFunc
	globalErrorWriter ErrorWriter = writeError
)

type Handler func(w http.ResponseWriter, r *http.Request) Writer
type Writer http.HandlerFunc

func IsHtmx(w http.ResponseWriter, r *http.Request) bool {
	return r.Header.Get(HX_REQUEST) == "true"
}

func noOpWriter(w http.ResponseWriter, r *http.Request) {}

func NoOp() Writer {
	return noOpWriter
}

func Component(template templ.Component) Writer {
	return func(w http.ResponseWriter, r *http.Request) {
		err := template.Render(r.Context(), w)
		if err != nil {
			globalErrorWriter(w, r, err)
			return
		}
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		return
	}
}

func Redirect(url string) Writer {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running redirect handler")
		if IsHtmx(w, r) {
			w.Header().Add(HX_REDIRECT, url)
			w.WriteHeader(http.StatusSeeOther)
		} else {
			w.Header().Add("Location", url)
			w.WriteHeader(http.StatusSeeOther)
		}
	}
}

func ErrorMsg(message string, status int) Writer {
	return noOpWriter
}

func ErrorComponent(component templ.Component, status int) Writer {
	return noOpWriter
}

func String(msg string, status int) Writer {
	return func(w http.ResponseWriter, r *http.Request) {
		n, err := w.Write([]byte(msg))
		if err != nil {
			globalErrorWriter(w, r, err)
		}
		w.Header().Add("Content-Length", strconv.Itoa(n))
		w.Header().Add("Content-Type", "text/plain")
	}
}

func Error(err error, status int) Writer {
	return noOpWriter
}

type BoundHandler[T any] func(http.ResponseWriter, *http.Request, *T) Writer

type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc
type ErrorWriter func(w http.ResponseWriter, r *http.Request, err error)

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch err.(type) {
	case *bindError:
		http.Error(w, err.(*bindError).msg, http.StatusBadRequest)
	default:
		http.Error(w, fmt.Sprintf("Handler error: %s", err.Error()), http.StatusInternalServerError)
	}
}

func applyMiddleware(f http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}

func UseGlobal(middlewares ...MiddlewareFunc) {
	globalMiddlewares = append(globalMiddlewares, middlewares...)
}
