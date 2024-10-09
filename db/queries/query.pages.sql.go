// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.pages.sql

package queries

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createPageWithPreview = `-- name: CreatePageWithPreview :one
INSERT INTO pages (id, user_id, url, title, description) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, url, title, description, image_url, created, updated
`

type CreatePageWithPreviewParams struct {
	ID          string
	UserID      uuid.UUID
	Url         string
	Title       pgtype.Text
	Description pgtype.Text
}

func (q *Queries) CreatePageWithPreview(ctx context.Context, arg CreatePageWithPreviewParams) (Page, error) {
	row := q.db.QueryRow(ctx, createPageWithPreview,
		arg.ID,
		arg.UserID,
		arg.Url,
		arg.Title,
		arg.Description,
	)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.ImageUrl,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const deletePageById = `-- name: DeletePageById :execrows
DELETE FROM pages 
WHERE id = $1
`

func (q *Queries) DeletePageById(ctx context.Context, id string) (int64, error) {
	result, err := q.db.Exec(ctx, deletePageById, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deletePageForUser = `-- name: DeletePageForUser :execrows
DELETE FROM pages
WHERE id = $1
AND user_id = $2
`

type DeletePageForUserParams struct {
	ID     string
	UserID uuid.UUID
}

func (q *Queries) DeletePageForUser(ctx context.Context, arg DeletePageForUserParams) (int64, error) {
	result, err := q.db.Exec(ctx, deletePageForUser, arg.ID, arg.UserID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const readPageById = `-- name: ReadPageById :one
SELECT id, user_id, url, title, description, image_url, created, updated FROM pages
WHERE id = $1
LIMIT 1
`

func (q *Queries) ReadPageById(ctx context.Context, id string) (Page, error) {
	row := q.db.QueryRow(ctx, readPageById, id)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.ImageUrl,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const readPagesByUserBetween = `-- name: ReadPagesByUserBetween :many
SELECT id, user_id, url, title, description, image_url, created, updated FROM pages 
WHERE user_id = $1
AND created BETWEEN $2 AND $3
ORDER BY created DESC
`

type ReadPagesByUserBetweenParams struct {
	UserID    uuid.UUID
	Created   pgtype.Timestamp
	Created_2 pgtype.Timestamp
}

func (q *Queries) ReadPagesByUserBetween(ctx context.Context, arg ReadPagesByUserBetweenParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, readPagesByUserBetween, arg.UserID, arg.Created, arg.Created_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Url,
			&i.Title,
			&i.Description,
			&i.ImageUrl,
			&i.Created,
			&i.Updated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readPagesByUserId = `-- name: ReadPagesByUserId :many
SELECT id, user_id, url, title, description, image_url, created, updated FROM pages
WHERE user_id = $1
ORDER BY created DESC
LIMIT $2 OFFSET $3
`

type ReadPagesByUserIdParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) ReadPagesByUserId(ctx context.Context, arg ReadPagesByUserIdParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, readPagesByUserId, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Url,
			&i.Title,
			&i.Description,
			&i.ImageUrl,
			&i.Created,
			&i.Updated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePagePreview = `-- name: UpdatePagePreview :exec
UPDATE pages SET
    title = $1,
    description = $2,
    image_url = $3,
    updated = now()
WHERE id = $4
`

type UpdatePagePreviewParams struct {
	Title       pgtype.Text
	Description pgtype.Text
	ImageUrl    pgtype.Text
	ID          string
}

func (q *Queries) UpdatePagePreview(ctx context.Context, arg UpdatePagePreviewParams) error {
	_, err := q.db.Exec(ctx, updatePagePreview,
		arg.Title,
		arg.Description,
		arg.ImageUrl,
		arg.ID,
	)
	return err
}