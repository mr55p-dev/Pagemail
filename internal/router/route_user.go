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
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/render/views"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/pkg/request"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

func loginRedirectUrl(router *Router) string {
	// generate the reset url
	urlAddr := url.URL{
		Scheme: router.proto,
		Host:   router.host,
		Path:   "login/google",
	}
	return urlAddr.String()
}

func (router *Router) GetLogin(w http.ResponseWriter, r *http.Request) {
	urlAddr := loginRedirectUrl(router)
	response.Component(views.Login(router.googleClientId, urlAddr), w, r)
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

	// Validate user
	user, err := auth.LoginPm(r.Context(), router.db, req.Email, []byte(req.Password))
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to authenticate user")
		response.Error(w, r, err)
		return
	}
	logger.InfoCtx(r.Context(), "User authenticated", "user-id", user.ID)

	// Write the user session
	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	auth.SetId(sess, user.ID)
	err = sess.Save(r, w)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to save session")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}

	logger.InfoCtx(r.Context(), "Written session")
	response.Redirect(w, r, URI_DASH)
	return
}

type PostLoginGoogleParams struct {
	Credential string `form:"credential"`
	CRSFToken  string `form:"g_csrf_token"`
}

func (router *Router) PostLoginGoogle(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	logger.InfoCtx(r.Context(), "Received google login request")
	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	req := request.BindRequest[PostLoginGoogleParams](w, r)
	if req == nil {
		return
	}

	// validate CSRF token
	if ok := auth.CSRFCheck(r, "g_csrf_token"); !ok {
		logger.ErrorCtx(r.Context(), "Failed to validate CSRF token")
		response.Error(w, r, pmerror.ErrCSRF)
		return
	}

	// decode the JWT
	logger.InfoCtx(r.Context(), "Got id token")
	token, err := auth.ValidateIdToken(r.Context(), req.Credential, router.googleClientId)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to validate id token")
		response.Error(w, r, err)
		return
	}

	// try to load the user. If not allowed, redirect to link account
	user, err := auth.HandleIdpRequest(r.Context(), router.db, token.Email, []byte(token.Subject))
	if err != nil {
		if errors.Is(err, pmerror.ErrMismatchAcc) {
			logger.InfoCtx(r.Context(), "User has not linked google account")
			sess.Values[SESS_G_TOKEN] = token
			err = sess.Save(r, w)
			if err != nil {
				logger.WithError(err).ErrorCtx(r.Context(), "Failed to save session")
				response.Error(w, r, pmerror.ErrUnspecified)
				return
			}
			response.Redirect(w, r, "/login/link")
			return
		}
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to auth with google")
		response.Error(w, r, err)
		return
	}

	// set the session
	auth.SetId(sess, user.ID)
	err = sess.Save(r, w)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to save session")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}
	response.Redirect(w, r, URI_DASH)
}

func (router *Router) GetLinkAccount(w http.ResponseWriter, r *http.Request) {
	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	token, ok := sess.Values[SESS_G_TOKEN].(auth.G_Token)
	if !ok {
		logger.ErrorCtx(r.Context(), "Could not load google cred from sesion")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}
	logger.InfoCtx(r.Context(), "Loaded google cred from session", "linking-to", token.Email)
	response.Component(render.LinkAccount(token.Email), w, r)
	return
}

type PostLinkAccountParams struct {
	Password string `form:"password"`
}

func (router *Router) PostLinkAccount(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := request.BindRequest[PostLoginRequest](w, r)
	if req == nil {
		return
	}
	logger.InfoCtx(r.Context(), "Received link account request")

	// Load the session
	sess, _ := router.Sessions.Get(r, auth.SessionKey)
	token, ok := sess.Values[SESS_G_TOKEN].(auth.G_Token)
	if !ok {
		logger.InfoCtx(r.Context(), "No google credential in session")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}

	// Login flow for pagemail platform
	user, err := auth.LoginPm(r.Context(), router.db, token.Email, []byte(req.Password))
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to authenticate user")
		response.Error(w, r, err)
		return
	}
	logger.InfoCtx(r.Context(), "user auth succesfully when linking google account")

	// Link the google account
	err = auth.LinkGoogleAccount(
		r.Context(),
		router.db,
		user,
		[]byte(token.Subject),
	)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to link google account")
		response.Error(w, r, err)
		return
	}
	logger.InfoCtx(r.Context(), "Linked google account")

	// set the session
	auth.SetId(sess, user.ID)
	delete(sess.Values, "google_subject")
	err = sess.Save(r, w)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to save session")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}
	response.Redirect(w, r, URI_DASH)
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
	passwordHash := auth.HashPassword([]byte(req.Password))
	_, tokenHash := auth.NewShortcutToken()
	user, err := queries.New(router.db).CreateUser(r.Context(), queries.CreateUserParams{
		ID:         tools.NewUserId(),
		Username:   req.Username,
		Email:      req.Email,
		Subscribed: req.Subscribed,
	})
	err = queries.New(router.db).CreateLocalAuth(r.Context(), queries.CreateLocalAuthParams{
		UserID:       user.ID,
		PasswordHash: passwordHash,
	})
	err = queries.New(router.db).CreateShortcutAuth(r.Context(), queries.CreateShortcutAuthParams{
		UserID:     user.ID,
		Credential: tokenHash,
	})
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

func (router *Router) GetSignup(w http.ResponseWriter, r *http.Request) {
	urlAddr := loginRedirectUrl(router)
	response.Component(views.Signup(router.googleClientId, urlAddr), w, r)

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

	user, err := queries.New(router.db).ReadUserByEmail(r.Context(), req.Email)
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

	err = queries.New(router.db).UpdateResetToken(r.Context(), queries.UpdateResetTokenParams{
		UserID:              user.ID,
		PasswordResetToken:  tokenHash,
		PasswordResetExpiry: sql.NullTime{Valid: true, Time: expires},
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Seting password reset token")
		response.Error(w, r, pmerror.NewInternalError("Failed to generate a reset token."))
	}

	// generate the reset url
	urlAddr := url.URL{
		Scheme: router.proto,
		Host:   router.host,
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
	userID, err := queries.New(router.db).ReadByResetToken(r.Context(), queries.ReadByResetTokenParams{
		PasswordResetToken:  hashedToken,
		PasswordResetExpiry: sql.NullTime{Valid: true, Time: now},
	})
	if err != nil {
		logger.WithError(err).Info("No user with valid reset token found", "token-hash", hashedToken)
		response.Error(w, r, pmerror.NewInternalError("Failed to update password"))
		return
	}

	// generate and update the password for the matching reset token
	n, err := queries.New(router.db).UpdatePassword(r.Context(), queries.UpdatePasswordParams{
		PasswordHash: hashedPassword,
		UserID:       userID,
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
	queries.New(router.db).UpdateResetToken(r.Context(), queries.UpdateResetTokenParams{
		PasswordResetToken:  []byte{},
		PasswordResetExpiry: sql.NullTime{},
		UserID:              userID,
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
