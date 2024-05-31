// //go:build integration

package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/middlewares"
	"github.com/stretchr/testify/assert"
)

type NilPreviewer struct{}

func (NilPreviewer) Queue(string) {}

var mux http.Handler
var session_cookie *http.Cookie

func init() {
	ctx := context.TODO()

	// setup the database
	conn := db.MustConnect(ctx, ":memory:")
	db.MustLoadSchema(ctx, conn)

	router, err := New(
		ctx,
		conn,
		nil,
		nil,
		&NilPreviewer{},
		strings.NewReader("passwordpassword"),
		"google_client",
		"", "", nil,
	)
	if err != nil {
		panic(err)
	}

	fn := middlewares.GetUserLoader(router.Sessions, router.db)
	mux = fn(router.Mux)
}

func WithHtmx(r *http.Request) *http.Request {
	r.Header.Add("Hx-Request", "true")
	return r
}

func TestRoot(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	res := w.Result()
	assert.Equal(res.StatusCode, http.StatusOK)
}

func TestSignup(t *testing.T) {
	assert := assert.New(t)
	form := strings.NewReader(url.Values(map[string][]string{
		"username":        {"test"},
		"email":           {"test@mail.com"},
		"password":        {"password"},
		"password-repeat": {"password"},
	}).Encode())
	r := httptest.NewRequest(http.MethodPost, "/signup", form)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, WithHtmx(r))

	res := w.Result()
	assert.Equal(res.StatusCode, http.StatusOK)

	cookies := res.Cookies()
	assert.Len(cookies, 1)
	cookie := cookies[0]
	assert.Equal(cookie.Name, auth.SessionKey)
	assert.NotZero(cookie.Value)
	assert.Greater(cookie.MaxAge, 0)
	assert.Equal(res.Header.Get("HX-Redirect"), "/pages/dashboard")
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	form := strings.NewReader(url.Values(map[string][]string{
		"email":    {"test@mail.com"},
		"password": {"password"},
	}).Encode())
	r := httptest.NewRequest(http.MethodPost, "/login/", form)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, WithHtmx(r))

	res := w.Result()
	assert.Equal(res.StatusCode, http.StatusOK)

	cookies := res.Cookies()
	assert.Len(cookies, 1)
	cookie := cookies[0]
	assert.Equal(cookie.Name, auth.SessionKey)
	assert.NotZero(cookie.Value)
	assert.Greater(cookie.MaxAge, 0)
	assert.Equal(res.Header.Get("HX-Redirect"), "/pages/dashboard")
	session_cookie = cookies[0]
}

func TestPostPage(t *testing.T) {
	assert := assert.New(t)
	form := strings.NewReader(url.Values(map[string][]string{
		"url": {"https://google.com"},
	}).Encode())
	r := httptest.NewRequest(http.MethodPost, "/pages/", form)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(session_cookie)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, WithHtmx(r))

	res := w.Result()
	assert.Equal(res.StatusCode, http.StatusOK)
}

func TestLogout(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest(http.MethodGet, "/user/logout", nil)
	r.AddCookie(session_cookie)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, WithHtmx(r))
	res := w.Result()
	assert.Equal(res.StatusCode, http.StatusOK)
}
