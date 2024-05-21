package auth

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

func commitOrRollback(ctx context.Context, tx *sql.Tx, err error) {
	if err != nil {
		logger.WithError(err).ErrorCtx(ctx, "Attempting to rollback transaction")
		txErr := tx.Rollback()
		if txErr != nil {
			logger.WithError(txErr).ErrorCtx(ctx, "Failed to rollback transaction")
		} else {
			logger.InfoCtx(ctx, "Rolled back user transaction")
		}
	} else {
		txErr := tx.Commit()
		if txErr != nil {
			logger.WithError(err).ErrorCtx(ctx, "Failed to commit transaction")
		} else {
			logger.InfoCtx(ctx, "Comitted user transaction")
		}
	}
}

func SignupUserIdp(ctx context.Context, db *sql.DB, email, name string, secret []byte) (queries.User, error) {
	var user queries.User
	var err error

	logger.InfoCtx(ctx, "Creating google IDP user")
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return user, err
	}
	defer commitOrRollback(ctx, tx, err)

	user, err = queries.New(db).WithTx(tx).CreateUser(ctx, queries.CreateUserParams{
		ID:         tools.GenerateNewId(10),
		Username:   name,
		Email:      email,
		Subscribed: true,
	})
	err = queries.New(db).WithTx(tx).CreateIdpAuth(ctx, queries.CreateIdpAuthParams{
		UserID:     user.ID,
		Platform:   "google",
		Credential: secret,
	})
	_, hashedToken := NewShortcutToken()
	err = queries.New(db).WithTx(tx).CreateShortcutAuth(ctx, queries.CreateShortcutAuthParams{
		UserID:     user.ID,
		Credential: hashedToken,
	})
	return user, err
}

func HandleIdpRequest(ctx context.Context, db *sql.DB, email string, cred []byte) (*queries.User, error) {
	// lookup the user by email
	// 1. If the user is new, then create a new record in the DB with their
	//	  info and log them in silently
	// 2. If the user is returning and has not previously auth'd with google
	//    then require them to login with their password to link their google account
	// 3. If the user is returning with google then just take them through
	q := queries.New(db)
	user, err := q.ReadUserByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.WithError(err).ErrorCtx(ctx, "Failed to lookup user for idp login")
			return &user, pmerror.ErrUnspecified
		}

		// User does not exist - create a new one
		logger.InfoCtx(ctx, "Creating new sign in with google user")
		newUser, err := SignupUserIdp(
			ctx,
			db,
			email,
			"",
			cred,
		)
		if err != nil {
			logger.WithError(err).ErrorCtx(ctx, "Failed to sign up idp user")
			return nil, pmerror.ErrUnspecified
		}
		logger.InfoCtx(ctx, "Created new IDP user")
		return &newUser, nil
	}

	// User does exist - look them up
	authMethods, err := q.ReadAuthMethods(ctx, user.ID)
	hasGoogleLink := false
	for _, method := range authMethods {
		// previous google sign in
		if method.Platform == "google" &&
			subtle.ConstantTimeCompare(method.Credential, []byte(cred)) == 1 {
			hasGoogleLink = true
		}
	}
	if !hasGoogleLink {
		// They do not have google linked
		// TODO: implement account linking
		logger.ErrorCtx(ctx, "requested to auth with google to a standard account")
		return nil, pmerror.ErrMismatchAcc
	}
	return &user, nil
}
