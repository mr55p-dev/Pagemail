package router

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/pkg/request"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

func (router *Router) GetLogin(w http.ResponseWriter, r *http.Request) {
	response.Component(render.Login(), w, r)
}

type PostLoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (router *Router) PostLogin(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := request.BindRequest[PostLoginRequest](w, r)
	if req == nil {
		response.Error(w, r, pmerror.ErrNoParam)
		return
	}

	logger.DebugCtx(r.Context(), "Received bound data", "email", req.Email, "req", req)
	user, err := router.DBClient.ReadUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(w, r, pmerror.ErrBadEmail)
		} else {
			response.Error(w, r, pmerror.NewInternalError("Sign-up failed"))
		}
		return
	}

	// Validate user
	if ok := auth.ValidateEmail([]byte(req.Email), []byte(user.Email)); !ok {
		logger.DebugCtx(r.Context(), "Invalid username")
		response.Error(w, r, pmerror.ErrBadEmail)
		return
	}
	if ok := auth.ValidatePassword([]byte(req.Password), user.Password); !ok {
		logger.DebugCtx(r.Context(), "Invalid password")
		response.Error(w, r, pmerror.ErrBadPassword)
		return
	}

	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	auth.SetId(sess, user.ID)
	sess.Save(r, w)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
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
	req := request.BindRequest[PostSignupRequest](w, r)
	if req == nil {
		return
	}

	// Read the form requests
	if req.Password != req.PasswordRepeat {
		response.Error(w, r, pmerror.ErrDiffPasswords)
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
			response.Error(w, r, pmerror.ErrDuplicateEmail)
			return
		}
	unknown:
		response.Error(w, r, pmerror.NewInternalError("Something went wrong signing you up"))
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
	response.Component(render.Signup(), w, r)
}

func (router *Router) GetPasswordReset(w http.ResponseWriter, r *http.Request) {
	response.Component(render.PasswordResetSubmit(), w, r)
}

func (router *Router) GetLogout(w http.ResponseWriter, r *http.Request) {
	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	auth.DelId(sess)
	_ = sess.Save(r, w)
	w.Header().Add("HX-Redirect", "/")
	return
}
