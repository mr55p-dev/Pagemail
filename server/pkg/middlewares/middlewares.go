package middlewares

import (
	"log/slog"
	"net/http"
	"path"

	"github.com/mr55p-dev/go-httpit"
	"github.com/mr55p-dev/go-httpit/pkg/trace"
)

func RequestLogger(log *slog.Logger, excludePaths ...string) httpit.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, v := range excludePaths {
				ok, _ := path.Match(v, r.URL.Path)
				if ok {
					goto next
				}
			}
			log.InfoContext(r.Context(), "Request", "method", r.Method, "path", r.URL.Path)
		next:
			next(w, r)
		}
	}
}

func Recover(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "Panic recover", "error", rec)
			}
		}()
		next(w, r)
	}
}

func Trace(idfunc func() string) httpit.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			id := idfunc()
			ctx := trace.SetTrace(r.Context(), id)
			r = r.WithContext(ctx)
			w.Header().Set(trace.TraceHeader, id)
			next(w, r)
		}
	}
}
