package router

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/a-h/templ"
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
	response.Redirect(w, r, "/pages/dashboard")
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
	if err := auth.CheckSubmittedPasswords([]byte(req.Password), []byte(req.PasswordRepeat)); err != nil {
		response.Error(w, r, err)
		return
	}

	// Generate a new user
	now := time.Now()
	passwordHash := auth.HashPassword([]byte(req.Password))
	_, tokenHash := auth.NewShortcutToken()
	user := dbqueries.CreateUserParams{
		ID:             tools.GenerateNewId(10),
		Username:       req.Username,
		Email:          req.Email,
		Password:       passwordHash,
		Subscribed:     req.Subscribed,
		ShortcutToken:  tokenHash,
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
	response.Redirect(w, r, "/pages/dashboard")
	return
}

func (Router) GetSignup(w http.ResponseWriter, r *http.Request) {
	response.Component(render.Signup(), w, r)
}

func (Router) GetPassResetReq(w http.ResponseWriter, r *http.Request) {
	response.Component(render.PasswordResetReq(), w, r)
}

type PostPasswordResetParams struct {
	Email string `form:"email"`
}

func (router *Router) PostPassResetReq(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := request.BindRequest[PostPasswordResetParams](w, r)
	if req == nil {
		return
	}

	user, err := router.DBClient.ReadUserByEmail(r.Context(), req.Email)
	if err != nil {
		logger.WithError(err).InfoCtx(r.Context(), "Error fetching user from DB")
		if errors.Is(err, sql.ErrNoRows) {
			response.Success("If such an email exists, check the inbox for a reset url.", w, r)
			return
		}
		logger.WithError(err).InfoCtx(r.Context(), "Error fetching user from DB")
		response.Error(w, r, pmerror.ErrBadEmail)
		return
	}

	// generate a new token and hash
	token, tokenHash := auth.NewResetToken()
	expires := time.Now().Add(time.Hour)

	err = router.DBClient.UpdateUserResetToken(r.Context(), dbqueries.UpdateUserResetTokenParams{
		ID:            user.ID,
		ResetToken:    tokenHash,
		ResetTokenExp: sql.NullTime{Valid: true, Time: expires},
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Seting password reset token")
		response.Error(w, r, pmerror.NewInternalError("Failed to generate a reset token."))
	}

	// generate the reset url
	urlAddr := url.URL{
		Scheme: "https",
		Host:   "pagemail.io",
		Path:   "password-reset/reset",
	}
	urlAddr.Query().Add("token", string(token))

	// Send an email
	buf := new(bytes.Buffer)
	err = render.PasswordResetMail(user.ID, templ.SafeURL(urlAddr.String())).Render(r.Context(), buf)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Writing password reset mail buffer")
		response.Error(w, r, pmerror.ErrCreatingMail)
		return
	}
	err = router.Sender.Send(r.Context(), user.Email, buf)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Sending password reset email")
		response.Error(w, r, pmerror.ErrCreatingMail)
		return
	}

	response.Success("Check your inbox for an email from support@pagemail.io", w, r)
}

func (Router) GetPassResetRedeem(w http.ResponseWriter, r *http.Request) {
	response.Component(render.PasswordReset(), w, r)
}

type PostPassResetRedeemParams struct {
	Token          string `query:"token"`
	Password       string `form:"password"`
	PasswordRepeat string `form:"password-repeat"`
}

func (router *Router) PostPassResetRedeem(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := request.BindRequest[PostPassResetRedeemParams](w, r)
	if req == nil {
		return
	}
	if err := auth.CheckSubmittedPasswords([]byte(req.Password), []byte(req.PasswordRepeat)); err != nil {
		response.Error(w, r, err)
		return
	}

	// generate and update the password for the matching reset token
	hashedPassword := auth.HashPassword([]byte(req.Password))
	err := router.DBClient.UpdateUserPassword(r.Context(), dbqueries.UpdateUserPasswordParams{
		Password:      hashedPassword,
		ResetToken:    []byte(req.Token),
		ResetTokenExp: sql.NullTime{Valid: true, Time: time.Now()},
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to update password")
		response.Error(w, r, pmerror.NewInternalError("Failed to update password"))
		return
	}

	logger.Info("Updated password")
	response.Success("Updated password", w, r)
}

func (router *Router) GetLogout(w http.ResponseWriter, r *http.Request) {
	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	auth.DelId(sess)
	_ = sess.Save(r, w)
	response.Redirect(w, r, "/")
	return
}
