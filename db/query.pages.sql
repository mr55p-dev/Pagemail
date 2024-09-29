-- name: CreatePageWithPreview :one
INSERT INTO pages (id, user_id, url, title, description) 
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ReadPageById :one
SELECT * FROM pages
WHERE id = $1
LIMIT 1;

-- name: UpdatePagePreview :exec
UPDATE pages SET
    title = $1,
    description = $2,
    image_url = $3,
    updated = now()
WHERE id = $4;

-- name: DeletePageById :execrows
DELETE FROM pages 
WHERE id = $1;

-- name: DeletePageForUser :execrows
DELETE FROM pages
WHERE id = $1
AND user_id = $2;

-- name: ReadPagesByUserId :many
SELECT * FROM pages
WHERE user_id = $1
ORDER BY created DESC
LIMIT $2 OFFSET $3;

-- name: ReadPagesByUserBetween :many
SELECT * FROM pages 
WHERE user_id = $1
AND created BETWEEN $2 AND $3
ORDER BY created DESC;
