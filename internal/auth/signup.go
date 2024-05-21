package auth

import (
	"context"
	"database/sql"

	"github.com/mr55p-dev/pagemail/db/queries"
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
