-- name: CreatePage :one
INSERT INTO pages (id, user_id, url, preview_state)
VALUES (?, ?, ?, 'unknown')
RETURNING *;

-- name: ReadPageById :one
SELECT * FROM pages
WHERE id = ?
LIMIT 1;

-- name: ReadPagesByReadable :many
SELECT * FROM pages
WHERE readable = ?
AND user_id = ?;

-- name: UpdatePagePreview :exec
UPDATE pages SET
    title = ?,
    description = ?,
    image_url = ?,
    preview_state = ?,
    updated = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdatePageReadability :exec
UPDATE pages SET
    readable = ?,
	updated = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeletePageById :execrows
DELETE FROM pages 
WHERE id = ?;

-- name: DeletePagesByUserId :execrows
DELETE FROM pages
WHERE user_id = ?;

-- name: ReadPagesByUserId :many
SELECT * FROM pages
WHERE user_id = ?
ORDER BY created DESC
LIMIT ? OFFSET ?;

-- name: ReadPagesByUserBetween :many
SELECT * FROM pages 
WHERE created BETWEEN sqlc.arg(start) AND sqlc.arg(end)
AND user_id = sqlc.arg(user_id)
ORDER BY created DESC;
