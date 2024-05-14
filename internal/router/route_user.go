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
	Subscribed     bool   `form:"subscribed"`
}

type PostLoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func GetLoginCookie(val string) *http.Cookie {
	return &http.Cookie{
		Name:     auth.SessionKey,
		Value:    val,
		Path:     "/",
		MaxAge:   864000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (Router) GetLogin(w http.ResponseWriter, r *http.Request) {
	componentRender(render.Login(), w, r)
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
	if err := auth.ValidateUser([]byte(req.Email), []byte(user.Email)); err != nil {
		logger.WithError(err).DebugCtx(r.Context(), "Error validating user")
		genericResponse(w, http.StatusUnauthorized)
		return
	}
	if err := auth.ValidatePassword([]byte(req.Password), user.Password); err != nil {
		logger.WithError(err).DebugCtx(r.Context(), "Error validating user")
		genericResponse(w, http.StatusUnauthorized)
		return
	}

	sess, _ := router.Authorizer.New(r, user.ID)
	_ = router.Authorizer.Save(r, w, sess)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}

func (router *Router) PostSignup(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
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
		Subscribed:     req.Subscribed,
		ShortcutToken:  tools.GenerateNewShortcutToken(),
		HasReadability: false,
		Created:        now,
		Updated:        now,
	}
	err := router.DBClient.CreateUser(r.Context(), user)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to create user")
		genericResponse(w, http.StatusInternalServerError)
		return
	}
	// Generate a token for the user from the session manager
	token, _ := router.Authorizer.New(r, user.ID)
	_ = router.Authorizer.Save(r, w, token)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}

func (Router) GetSignup(w http.ResponseWriter, r *http.Request) {
	componentRender(render.Signup(), w, r)
	return
}

func (Router) GetLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   auth.SessionKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/pages/dashboard")
	return
}
