// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.auth.sql

package queries

import (
	"context"
	"database/sql"
)

const createIdpAuth = `-- name: CreateIdpAuth :exec
INSERT INTO auth (
    user_id,
    platform,
    credential
) VALUES (?, ?, ?)
`

type CreateIdpAuthParams struct {
	UserID     string
	Platform   string
	Credential []byte
}

func (q *Queries) CreateIdpAuth(ctx context.Context, arg CreateIdpAuthParams) error {
	_, err := q.db.ExecContext(ctx, createIdpAuth, arg.UserID, arg.Platform, arg.Credential)
	return err
}

const createLocalAuth = `-- name: CreateLocalAuth :exec
INSERT INTO auth (
    user_id,
    platform,
    password_hash
) VALUES (?, 'pagemail', ?)
`

type CreateLocalAuthParams struct {
	UserID       string
	PasswordHash []byte
}

func (q *Queries) CreateLocalAuth(ctx context.Context, arg CreateLocalAuthParams) error {
	_, err := q.db.ExecContext(ctx, createLocalAuth, arg.UserID, arg.PasswordHash)
	return err
}

const createShortcutAuth = `-- name: CreateShortcutAuth :exec
INSERT INTO auth (
    user_id,
    platform,
    credential
) VALUES (?, 'shortcut', ?)
`

type CreateShortcutAuthParams struct {
	UserID     string
	Credential []byte
}

func (q *Queries) CreateShortcutAuth(ctx context.Context, arg CreateShortcutAuthParams) error {
	_, err := q.db.ExecContext(ctx, createShortcutAuth, arg.UserID, arg.Credential)
	return err
}

const readAuthMethods = `-- name: ReadAuthMethods :many
SELECT id, user_id, platform, password_hash, password_reset_token, password_reset_expiry, credential, created, updated FROM auth WHERE user_id = ?
`

func (q *Queries) ReadAuthMethods(ctx context.Context, userID string) ([]Auth, error) {
	rows, err := q.db.QueryContext(ctx, readAuthMethods, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Auth
	for rows.Next() {
		var i Auth
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Platform,
			&i.PasswordHash,
			&i.PasswordResetToken,
			&i.PasswordResetExpiry,
			&i.Credential,
			&i.Created,
			&i.Updated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readByResetToken = `-- name: ReadByResetToken :one
SELECT user_id
FROM auth
WHERE password_reset_token = ?
    AND password_reset_expiry > ?
LIMIT 1
`

type ReadByResetTokenParams struct {
	PasswordResetToken  []byte
	PasswordResetExpiry sql.NullTime
}

func (q *Queries) ReadByResetToken(ctx context.Context, arg ReadByResetTokenParams) (string, error) {
	row := q.db.QueryRowContext(ctx, readByResetToken, arg.PasswordResetToken, arg.PasswordResetExpiry)
	var user_id string
	err := row.Scan(&user_id)
	return user_id, err
}

const readByUidPlatform = `-- name: ReadByUidPlatform :one
SELECT id, user_id, platform, password_hash, password_reset_token, password_reset_expiry, credential, created, updated
FROM auth
WHERE user_id = ?
AND platform = ?
LIMIT 1
`

type ReadByUidPlatformParams struct {
	UserID   string
	Platform string
}

func (q *Queries) ReadByUidPlatform(ctx context.Context, arg ReadByUidPlatformParams) (Auth, error) {
	row := q.db.QueryRowContext(ctx, readByUidPlatform, arg.UserID, arg.Platform)
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

const readUserByShortcut = `-- name: ReadUserByShortcut :one
SELECT users.id, users.email, users.username, users.subscribed, users.has_readability, users.created, users.updated
FROM users
LEFT JOIN auth
    ON auth.user_id = users.id
    AND auth.platform = 'shortcut'
WHERE auth.credential = ?
LIMIT 1
`

func (q *Queries) ReadUserByShortcut(ctx context.Context, credential []byte) (User, error) {
	row := q.db.QueryRowContext(ctx, readUserByShortcut, credential)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Subscribed,
		&i.HasReadability,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const updatePassword = `-- name: UpdatePassword :execrows
UPDATE auth
SET password_hash = ?
WHERE user_id = ?
    AND platform = 'pagemail'
`

type UpdatePasswordParams struct {
	PasswordHash []byte
	UserID       string
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, updatePassword, arg.PasswordHash, arg.UserID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateResetToken = `-- name: UpdateResetToken :exec
UPDATE auth
SET 
    password_reset_token = ?,
    password_reset_expiry = ?
WHERE user_id = ?
    AND platform = 'pagemail'
`

type UpdateResetTokenParams struct {
	PasswordResetToken  []byte
	PasswordResetExpiry sql.NullTime
	UserID              string
}

func (q *Queries) UpdateResetToken(ctx context.Context, arg UpdateResetTokenParams) error {
	_, err := q.db.ExecContext(ctx, updateResetToken, arg.PasswordResetToken, arg.PasswordResetExpiry, arg.UserID)
	return err
}

const updateShortcutToken = `-- name: UpdateShortcutToken :exec
UPDATE auth
SET credential = ?
WHERE user_id = ?
    AND platform = 'shortcut'
`

type UpdateShortcutTokenParams struct {
	Credential []byte
	UserID     string
}

func (q *Queries) UpdateShortcutToken(ctx context.Context, arg UpdateShortcutTokenParams) error {
	_, err := q.db.ExecContext(ctx, updateShortcutToken, arg.Credential, arg.UserID)
	return err
}