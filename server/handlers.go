package httpit

import (
	"net/http"
)

func NewMiddlewareHandler(fn http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	fn = applyMiddleware(fn, middlewares...)
	fn = applyMiddleware(fn, globalMiddlewares...)
	return fn
}

func NewHandler(handler Handler, middlewares ...MiddlewareFunc) http.HandlerFunc {
	return NewMiddlewareHandler(func(w http.ResponseWriter, r *http.Request) {
		ctx := newContext(w, r)
		err := handler(ctx)
		if err != nil {
			WriteError(w, err)
			return
		}

		ctx.Finalize(w)
		w.WriteHeader(http.StatusNoContent)
		return

	})
}

func NewInHandler[T any](handler InHandler[T], middlewares ...MiddlewareFunc) http.HandlerFunc {
	return NewMiddlewareHandler(func(w http.ResponseWriter, r *http.Request) {
		in := new(T)
		if err := bind(in, r); err != nil {
			http.Error(w, "Failed to bind request parameters", http.StatusBadRequest)
		}

		ctx := newContext(w, r)
		err := handler(ctx, in)
		if err != nil {
			WriteError(w, err)
			return
		}

		ctx.Finalize(w)
		w.WriteHeader(http.StatusNoContent)
		return

	})
}

func NewTemplHandler(handler TemplHandler, middlewares ...MiddlewareFunc) http.HandlerFunc {
	return NewMiddlewareHandler(func(w http.ResponseWriter, r *http.Request) {
		ctx := newContext(w, r)
		component, err := handler(ctx)
		if err != nil {
			WriteError(w, err)
			return
		}
		err = component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			return
		}
		ctx.Finalize(w)
		w.Header().Set("Content-Type", "text/html")
		return

	})
}

func NewTemplMappedHandler[In any](handler TemplMappedHandler[In], middlewares ...MiddlewareFunc) http.HandlerFunc {
	return NewMiddlewareHandler(func(w http.ResponseWriter, r *http.Request) {
		inp := new(In)
		if err := bind(inp, r); err != nil {
			http.Error(w, "Failed to bind request parameters", http.StatusBadRequest)
			return
		}

		ctx := newContext(w, r)
		component, err := handler(ctx, inp)
		if err != nil {
			WriteError(w, err)
			return
		}
		err = component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		return

	})
}
