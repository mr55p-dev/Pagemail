package httpit

import (
	"net/http"
)

func withMiddleware(fn http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	fn = applyMiddleware(fn, middlewares...)
	fn = applyMiddleware(fn, globalMiddlewares...)
	return fn
}

func NewHandler(handler Handler, middlewares ...MiddlewareFunc) http.HandlerFunc {
	return withMiddleware(func(w http.ResponseWriter, r *http.Request) {
		mod := handler(w, r)
		if mod != nil {
			mod(w, r)
		}
	}, middlewares...)
}

func NewBoundHandler[T any](handler BoundHandler[T], middlewares ...MiddlewareFunc) http.HandlerFunc {
	return withMiddleware(func(w http.ResponseWriter, r *http.Request) {
		in := new(T)
		if err := bind(in, r); err != nil {
			globalErrorWriter(w, r, err)
		}

		mod := handler(w, r, in)
		if mod != nil {
			mod(w, r)
		}
	}, middlewares...)
}
