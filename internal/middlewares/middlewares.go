package middlewares

import (
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/internal/trace"
)

type MiddlewareFunc func(next http.Handler) http.Handler

func RequestLogger(next http.Handler) http.Handler {
	log := logging.NewLogger("request-logger")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.InfoCtx(r.Context(), "Request",
			"method", r.Method,
			"url", r.URL.Path,
			"remote", r.RemoteAddr,
			"trace-id", trace.GetTraceId(r.Context()),
		)
		next.ServeHTTP(w, r)
	})
}

func Tracer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := tools.GenerateNewId(10)
		ctx := trace.SetTrace(r.Context(), traceId)
		w.Header().Add("X-Trace-Id", traceId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Recover(next http.Handler) http.Handler {
	logger := logging.NewLogger("middleware-recover")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.WithRecover(rec).ErrorCtx(r.Context(), "Recovered from panic")

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func WithMiddleware(mux http.Handler, middleware ...MiddlewareFunc) http.Handler {
	for i, j := 0, len(middleware)-1; i < j; i, j = i+1, j-1 {
		middleware[i], middleware[j] = middleware[j], middleware[i]
	}
	for _, m := range middleware {
		mux = m(mux)
	}

	return mux
}
