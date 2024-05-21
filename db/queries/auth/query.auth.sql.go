// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.auth.sql

package queryauth

import (
	"context"
)

const readUserByResetToken = `-- name: ReadUserByResetToken :one
SELECT id, user_id, platform, password_hash, password_reset_token, password_reset_expiry, credential, created, updated FROM auth LIMIT 1
`

func (q *Queries) ReadUserByResetToken(ctx context.Context) (Auth, error) {
	row := q.db.QueryRowContext(ctx, readUserByResetToken)
	var i Auth
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Platform,
		&i.PasswordHash,
		&i.PasswordResetToken,
		&i.PasswordResetExpiry,
		&i.Credential,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :execrows
UPDATE auth
SET password_hash = ?
WHERE user_id = ?
    AND platform = 'pagemail'
RETURNING user_id
`

type UpdateUserPasswordParams struct {
	PasswordHash []byte
	UserID       string
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, updateUserPassword, arg.PasswordHash, arg.UserID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
