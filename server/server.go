package httpit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

type Context interface {
	context.Context
	SetRedirect(string)
	SetCookie(http.Cookie)
	IsHtmx() bool
}

type reqCtx struct {
	context.Context
	r           *http.Request
	w           http.ResponseWriter
	redirectUrl string
	cookies     []http.Cookie
}

func (ctx *reqCtx) SetRedirect(url string) {
	ctx.redirectUrl = url
}

func (ctx *reqCtx) SetCookie(cookie http.Cookie) {
	ctx.cookies = append(ctx.cookies, cookie)
}

func (ctx *reqCtx) IsHtmx() bool {
	return ctx.r.Header.Get("HX-Request") != ""
}

func (ctx *reqCtx) Finalize(w http.ResponseWriter) {
	for _, v := range ctx.cookies {
		http.SetCookie(w, &v)
	}
	if ctx.redirectUrl != "" {
		w.Header().Add("HX-Redirect", ctx.redirectUrl)
	}
}

func newContext(w http.ResponseWriter, r *http.Request) *reqCtx {
	return &reqCtx{
		Context: r.Context(),
		r:       r,
		w:       w,
	}
}

type Handler func(Context) error
type InHandler[T any] func(Context, *T) error
type TemplHandler func(Context) (templ.Component, error)
type TemplMappedHandler[TIn any] func(Context, *TIn) (templ.Component, error)

type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc
type ErrWriter func() HttpError

type HttpError struct {
	Code int
	Msg  string
}

var globalMiddlewares []MiddlewareFunc

func (e *HttpError) Error() string {
	return fmt.Sprintf("Handler error (%d): %s", e.Code, e.Msg)
}

func Error(msg string, code int) *HttpError {
	return &HttpError{
		Msg:  msg,
		Code: code,
	}
}

func WriteError(w http.ResponseWriter, err error) {
	httpErr, ok := err.(*HttpError)
	if ok {
		w.WriteHeader(httpErr.Code)
		w.Write([]byte(httpErr.Msg))
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func bind[T any](in T, r *http.Request) error {
	return nil
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
