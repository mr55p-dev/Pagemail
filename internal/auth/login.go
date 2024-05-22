package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
)

func CSRFCheck(r *http.Request, name string) bool {
	cookie, err := r.Cookie(name)
	if err != nil {
		logger.WithError(err).InfoCtx(r.Context(), "failed to load CSRF cookie")
		return false
	}
	if cookie.Value == "" {
		logger.InfoCtx(r.Context(), "No body in CSRF cookie")
		return false
	}
	param := r.FormValue(name)
	if param == "" {
		logger.InfoCtx(r.Context(), "No body in CRSF form")
		return false
	}
	ok := param == cookie.Value
	if !ok {
		logger.WarnCtx(r.Context(), "CRSF token check failed")
	}
	return ok
}

func LoginIdp(ctx context.Context, db *sql.DB) {
	
}

func LoginPm(ctx context.Context, db *sql.DB, email string, password []byte) (queries.User, error) {
	var user queries.User
	var err error

	// Read the user
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	defer commitOrRollback(ctx, tx, err)
	if err != nil {
		return user, pmerror.ErrUnspecified
	}
	user, err = queries.New(db).WithTx(tx).ReadUserByEmail(ctx, email)
	if err != nil {
		logger.WithError(err).ErrorCtx(ctx, "Failed to read user from DB")
		if errors.Is(err, sql.ErrNoRows) {
			return user, pmerror.ErrBadEmail
		} else {
			return user, pmerror.ErrUnspecified
		}
	}
	logger.InfoCtx(ctx, "User found", "user-id", user.ID)
	authRecord, err := queries.New(db).WithTx(tx).ReadByUidPlatform(ctx, queries.ReadByUidPlatformParams{
		UserID:   user.ID,
		Platform: "pagemail",
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(ctx, "Failed to read auth for user")
		if errors.Is(err, sql.ErrNoRows) {
			return user, pmerror.ErrNoAuth
		} else {
			return user, pmerror.ErrUnspecified
		}
	}

	if ok := ValidateEmail([]byte(email), []byte(user.Email)); !ok {
		logger.InfoCtx(ctx, "Invalid username")
		return user, pmerror.ErrBadEmail
	}
	if ok := ValidatePassword([]byte(password), authRecord.PasswordHash); !ok {
		logger.InfoCtx(ctx, "Invalid password")
		return user, pmerror.ErrBadPassword
	}
	return user, nil
}
