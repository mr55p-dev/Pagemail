package router

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

type PostSignupRequest struct {
	Username       string `form:"username"`
	Email          string `form:"email"`
	Password       string `form:"password"`
	PasswordRepeat string `form:"password-repeat"`
}

type PostLoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func GetLoginCookie(val string) *http.Cookie {
	return &http.Cookie{
		Name:     auth.SESS_COOKIE,
		Value:    val,
		Path:     "/",
		MaxAge:   864000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (Router) GetLogin(w http.ResponseWriter, r *http.Request) {
	staticRender(render.Login(), w, r)
}

func (router *Router) PostLogin(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := requestBind[PostLoginRequest](w, r)
	if req == nil {
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
	err = auth.ValidateUser([]byte(req.Email), []byte(user.Email), []byte(req.Password), user.Password)
	if err != nil {
		// TODO: handle the different auth errors
		logger.WithError(err).DebugCtx(r.Context(), "Error validating user")
		genericResponse(w, http.StatusUnauthorized)
		return
	}

	sess := router.Authorizer.GenSessionToken(user.ID)
	cookie := GetLoginCookie(sess)

	http.SetCookie(w, cookie)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}

func (router *Router) PostSignup(w http.ResponseWriter, r *http.Request) {
	req := requestBind[PostSignupRequest](w, r)
	if req == nil {
		return
	}

	// Read the form requests
	if req.Password != req.PasswordRepeat {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
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
		Subscribed:     false,
		ShortcutToken:  tools.GenerateNewShortcutToken(),
		HasReadability: false,
		Created:        now,
		Updated:        now,
	}
	err := router.DBClient.CreateUser(r.Context(), user)
	if err != nil {
		genericResponse(w, http.StatusInternalServerError)
		return
	}

	// Generate a token for the user from the session manager
	token := router.Authorizer.GenSessionToken(user.ID)
	cookie := GetLoginCookie(token)

	http.SetCookie(w, cookie)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}

func (Router) GetSignup(w http.ResponseWriter, r *http.Request) {
	staticRender(render.Signup(), w, r)
	return
}

func (Router) GetLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   auth.SESS_COOKIE,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}