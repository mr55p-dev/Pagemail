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
	"github.com/mr55p-dev/pagemail/internal/mail"
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
	logger.InfoCtx(r.Context(), "Received login request")
	req := request.BindRequest[PostLoginRequest](w, r)
	if req == nil {
		logger.ErrorCtx(r.Context(), "Failed to bind request parameters")
		response.Error(w, r, pmerror.ErrNoParam)
		return
	}

	user, err := router.DBClient.ReadUserByEmail(r.Context(), req.Email)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to read user from DB")
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(w, r, pmerror.ErrBadEmail)
		} else {
			response.Error(w, r, pmerror.ErrUnspecified)
		}
		return
	}
	logger.InfoCtx(r.Context(), "User found", "user-id", user.ID)

	// Validate user
	if ok := auth.ValidateEmail([]byte(req.Email), []byte(user.Email)); !ok {
		logger.InfoCtx(r.Context(), "Invalid username")
		response.Error(w, r, pmerror.ErrBadEmail)
		return
	}
	if ok := auth.ValidatePassword([]byte(req.Password), user.Password); !ok {
		logger.InfoCtx(r.Context(), "Invalid password")
		response.Error(w, r, pmerror.ErrBadPassword)
		return
	}
	logger.InfoCtx(r.Context(), "User authenticated", "user-id", user.ID)

	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	auth.SetId(sess, user.ID)
	err = sess.Save(r, w)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to save session")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}

	logger.InfoCtx(r.Context(), "Written session")
	response.Redirect(w, r, "/pages/dashboard")
	return
}

func (router *Router) PostLoginGoogle(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	logger.InfoCtx(r.Context(), "Received google login request")
	// validate credential
	email := "hello@mail.com"

	// lookup the user by email
	// if the user does not exist, create a new user
	user, err := router.DBClient.ReadUserByEmail(r.Context(), email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.WithError(err).ErrorCtx(r.Context(), "Failed to read user from DB")
			response.Error(w, r, pmerror.ErrUnspecified)
			return
		}
		user, err = auth.SignupUserIdp(r.Context(), router.DBClient, email, "Google user")
		if err != nil {
			logger.WithError(err).ErrorCtx(r.Context(), "Failed to create user")
			response.Error(w, r, pmerror.ErrUnspecified)
			return
		}
	}

	// set the session
	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	auth.SetId(sess, user.ID)
	err = sess.Save(r, w)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to save session")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}
	response.Redirect(w, r, "/pages/dashboard")
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
		Path:   "password-reset/redeem",
	}
	q := urlAddr.Query()
	q.Add("token", string(token))
	urlAddr.RawQuery = q.Encode()

	// Send an email
	buf := new(bytes.Buffer)
	err = render.PasswordResetMail(user.ID, templ.SafeURL(urlAddr.String())).Render(r.Context(), buf)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Writing password reset mail buffer")
		response.Error(w, r, pmerror.ErrCreatingMail)
		return
	}
	msg := mail.MakeMessage(
		user.Email,
		mail.WithSender("support@pagemail.io"),
		mail.WithSubject("Reset your password"),
		mail.WithBody(buf),
	)
	err = router.Sender.Send(r.Context(), msg)
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
	Token          string `form:"token"`
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

	// get the user
	hashedPassword := auth.HashPassword([]byte(req.Password))
	hashedToken := auth.HashValue([]byte(req.Token))
	now := time.Now()
	user, err := router.DBClient.ReadUserByResetToken(r.Context(), dbqueries.ReadUserByResetTokenParams{
		ResetToken:    hashedToken,
		ResetTokenExp: sql.NullTime{Valid: true, Time: now},
	})
	if err != nil {
		logger.WithError(err).Info("No user with valid reset token found", "token-hash", hashedToken)
		response.Error(w, r, pmerror.NewInternalError("Failed to update password"))
		return
	}

	// generate and update the password for the matching reset token
	n, err := router.DBClient.UpdateUserPassword(r.Context(), dbqueries.UpdateUserPasswordParams{
		Password:      hashedPassword,
		ResetToken:    hashedToken,
		ResetTokenExp: sql.NullTime{Valid: true, Time: now},
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to update password")
		response.Error(w, r, pmerror.NewInternalError("Failed to update password"))
		return
	}
	if n == 0 {
		logger.ErrorCtx(r.Context(), "No rows affected when updating password")
		response.Error(w, r, pmerror.NewInternalError("Failed to update password"))
		return
	}

	// clear the password reset token
	router.DBClient.UpdateUserResetToken(r.Context(), dbqueries.UpdateUserResetTokenParams{
		ResetToken:    []byte{},
		ResetTokenExp: sql.NullTime{},
		ID:            user.ID,
	})

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
