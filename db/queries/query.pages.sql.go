// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.pages.sql

package queries

import (
	"context"
	"database/sql"
	"time"
)

const createPage = `-- name: CreatePage :one
INSERT INTO pages (id, user_id, url, preview_state)
VALUES (?, ?, ?, 'unknown')
RETURNING id, user_id, url, title, description, image_url, preview_state, created, updated
`

type CreatePageParams struct {
	ID     string
	UserID string
	Url    string
}

func (q *Queries) CreatePage(ctx context.Context, arg CreatePageParams) (Page, error) {
	row := q.db.QueryRowContext(ctx, createPage, arg.ID, arg.UserID, arg.Url)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.ImageUrl,
		&i.PreviewState,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const deletePageById = `-- name: DeletePageById :execrows
DELETE FROM pages 
WHERE id = ?
`

func (q *Queries) DeletePageById(ctx context.Context, id string) (int64, error) {
	result, err := q.db.ExecContext(ctx, deletePageById, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const deletePagesByUserId = `-- name: DeletePagesByUserId :execrows
DELETE FROM pages
WHERE user_id = ?
`

func (q *Queries) DeletePagesByUserId(ctx context.Context, userID string) (int64, error) {
	result, err := q.db.ExecContext(ctx, deletePagesByUserId, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const readPageById = `-- name: ReadPageById :one
SELECT id, user_id, url, title, description, image_url, preview_state, created, updated FROM pages
WHERE id = ?
LIMIT 1
`

func (q *Queries) ReadPageById(ctx context.Context, id string) (Page, error) {
	row := q.db.QueryRowContext(ctx, readPageById, id)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.ImageUrl,
		&i.PreviewState,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const readPagesByUserBetween = `-- name: ReadPagesByUserBetween :many
SELECT id, user_id, url, title, description, image_url, preview_state, created, updated FROM pages 
WHERE created BETWEEN ?1 AND ?2
AND user_id = ?3
ORDER BY created DESC
`

type ReadPagesByUserBetweenParams struct {
	Start  time.Time
	End    time.Time
	UserID string
}

func (q *Queries) ReadPagesByUserBetween(ctx context.Context, arg ReadPagesByUserBetweenParams) ([]Page, error) {
	rows, err := q.db.QueryContext(ctx, readPagesByUserBetween, arg.Start, arg.End, arg.UserID)
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
			&i.PreviewState,
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

const readPagesByUserId = `-- name: ReadPagesByUserId :many
SELECT id, user_id, url, title, description, image_url, preview_state, created, updated FROM pages
WHERE user_id = ?
ORDER BY created DESC
LIMIT ? OFFSET ?
`

type ReadPagesByUserIdParams struct {
	UserID string
	Limit  int64
	Offset int64
}

func (q *Queries) ReadPagesByUserId(ctx context.Context, arg ReadPagesByUserIdParams) ([]Page, error) {
	rows, err := q.db.QueryContext(ctx, readPagesByUserId, arg.UserID, arg.Limit, arg.Offset)
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
			&i.PreviewState,
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

const updatePagePreview = `-- name: UpdatePagePreview :exec
UPDATE pages SET
    title = ?,
    description = ?,
    image_url = ?,
    preview_state = ?,
    updated = CURRENT_TIMESTAMP
WHERE id = ?
`

type UpdatePagePreviewParams struct {
	Title        sql.NullString
	Description  sql.NullString
	ImageUrl     sql.NullString
	PreviewState string
	ID           string
}

func (q *Queries) UpdatePagePreview(ctx context.Context, arg UpdatePagePreviewParams) error {
	_, err := q.db.ExecContext(ctx, updatePagePreview,
		arg.Title,
		arg.Description,
		arg.ImageUrl,
		arg.PreviewState,
		arg.ID,
	)
	return err
}
