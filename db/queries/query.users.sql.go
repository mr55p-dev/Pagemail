// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.users.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    username
) VALUES ($1, $2, $3)
RETURNING id, email, username, has_readability, created, updated
`

type CreateUserParams struct {
	ID       pgtype.UUID
	Email    string
	Username string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.ID, arg.Email, arg.Username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.HasReadability,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const readUserByEmail = `-- name: ReadUserByEmail :one
SELECT id, email, username, has_readability, created, updated FROM users 
WHERE email = $1
LIMIT 1
`

func (q *Queries) ReadUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, readUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.HasReadability,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const readUserById = `-- name: ReadUserById :one
SELECT id, email, username, has_readability, created, updated FROM users 
WHERE id = $1 
LIMIT 1
`

func (q *Queries) ReadUserById(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, readUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.HasReadability,
		&i.Created,
		&i.Updated,
	)
	return i, err
}
