package router

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

func (router *Router) GetLogin(w http.ResponseWriter, r *http.Request) {
	componentRender(render.Login(), w, r)
}

type PostLoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (router *Router) PostLogin(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := requestBind[PostLoginRequest](w, r)
	if req == nil {
		genericResponse(w, http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Received bound data", "email", req.Email, "req", req)
	user, err := router.DBClient.ReadUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			genericResponse(w, http.StatusUnauthorized)
		} else {
			genericResponse(w, http.StatusInternalServerError)
		}
		return
	}

	// Validate user
	if ok := auth.ValidateEmail([]byte(req.Email), []byte(user.Email)); !ok {
		logger.DebugCtx(r.Context(), "Invalid username")
		genericResponse(w, http.StatusUnauthorized)
		return
	}
	if ok := auth.ValidatePassword([]byte(req.Password), user.Password); !ok {
		logger.DebugCtx(r.Context(), "Invalid password")
		genericResponse(w, http.StatusUnauthorized)
		return
	}

	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	auth.SetId(sess, user.ID)
	sess.Save(r, w)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}

func errorResponse(w http.ResponseWriter, r *http.Request, detail string, status int) {
	if isHtmx(r) {
		w.WriteHeader(status)
		componentRender(render.ErrorBox("Error", detail), w, r)
	} else {
		http.Error(w, fmt.Sprintf("Error: %s", detail), status)
	}
}

type PostSignupRequest struct {
	Username       string `form:"username"`
	Email          string `form:"email"`
	Password       string `form:"password"`
	PasswordRepeat string `form:"password-repeat"`
	Subscribed     bool   `form:"subscribed"`
}

func (router *Router) PostSignup(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := requestBind[PostSignupRequest](w, r)
	if req == nil {
		return
	}

	// Read the form requests
	if req.Password != req.PasswordRepeat {
		errorResponse(w, r, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Generate a new user
	now := time.Now()
	passwordHash := auth.HashPassword([]byte(req.Password))
	user := dbqueries.CreateUserParams{
		ID:             tools.GenerateNewId(10),
		Username:       req.Username,
		Email:          req.Email,
		Password:       passwordHash,
		Avatar:         sql.NullString{},
		Subscribed:     req.Subscribed,
		ShortcutToken:  tools.GenerateNewShortcutToken(),
		HasReadability: false,
		Created:        now,
		Updated:        now,
	}
	err := router.DBClient.CreateUser(r.Context(), user)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to create user")
		sqlErr, ok := err.(sqlite3.Error)
		if !ok {
			goto unknown
		}
		if sqlErr.Code == sqlite3.ErrConstraint {
			errorResponse(w, r, "Looks like that email address is already taken. If you can't remember your password please reach out to help@pagemail.io for assistence.", http.StatusBadRequest)
			return
		}
	unknown:
		errorResponse(w, r, "Something went wrong signing you up", http.StatusInternalServerError)
		return
	}
	// Generate a token for the user from the session manager
	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	auth.SetId(sess, user.ID)
	sess.Save(r, w)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}

func (Router) GetSignup(w http.ResponseWriter, r *http.Request) {
	componentRender(render.Signup(), w, r)
	return
}

func (router *Router) GetLogout(w http.ResponseWriter, r *http.Request) {
	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	auth.DelId(sess)
	_ = sess.Save(r, w)
	w.Header().Add("HX-Redirect", "/")
	return
}
