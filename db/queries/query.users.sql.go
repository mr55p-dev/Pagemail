// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.users.sql

package queries

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    username,
    subscribed
) VALUES (?, ?, ?, ?)
RETURNING id, email, username, subscribed, has_readability, created, updated
`

type CreateUserParams struct {
	ID         string
	Email      string
	Username   string
	Subscribed bool
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Username,
		arg.Subscribed,
	)
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

const readUserByEmail = `-- name: ReadUserByEmail :one
SELECT id, email, username, subscribed, has_readability, created, updated FROM users 
WHERE email = ?
LIMIT 1
`

func (q *Queries) ReadUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, readUserByEmail, email)
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

const readUserById = `-- name: ReadUserById :one
SELECT id, email, username, subscribed, has_readability, created, updated FROM users 
WHERE id = ? 
LIMIT 1
`

func (q *Queries) ReadUserById(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, readUserById, id)
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

const readUsersWithMail = `-- name: ReadUsersWithMail :many
SELECT id, username, email FROM users 
WHERE subscribed = true
`

type ReadUsersWithMailRow struct {
	ID       string
	Username string
	Email    string
}

func (q *Queries) ReadUsersWithMail(ctx context.Context) ([]ReadUsersWithMailRow, error) {
	rows, err := q.db.QueryContext(ctx, readUsersWithMail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadUsersWithMailRow
	for rows.Next() {
		var i ReadUsersWithMailRow
		if err := rows.Scan(&i.ID, &i.Username, &i.Email); err != nil {
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

const updateUserSubscription = `-- name: UpdateUserSubscription :exec
UPDATE users SET 
subscribed = ? 
WHERE id = ?
`

type UpdateUserSubscriptionParams struct {
	Subscribed bool
	ID         string
}

func (q *Queries) UpdateUserSubscription(ctx context.Context, arg UpdateUserSubscriptionParams) error {
	_, err := q.db.ExecContext(ctx, updateUserSubscription, arg.Subscribed, arg.ID)
	return err
}
