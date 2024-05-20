-- name: CreatePage :exec
INSERT INTO pages (id, user_id, url, preview_state, created, updated)
VALUES (?, ?, ?, 'unknown', ?, ?);

-- name: ReadPageById :one
SELECT * FROM pages
WHERE id = ?
LIMIT 1;

-- name: UpdatePagePreview :exec
UPDATE pages SET
    title = ?,
    description = ?,
    image_url = ?,
    preview_state = ?,
    updated = ?
WHERE id = ?;

-- name: DeletePageById :execrows
DELETE FROM pages 
WHERE id = ?;

-- name: ReadPagesByUserId :many
SELECT * FROM pages
WHERE user_id = ?
ORDER BY created DESC;

-- name: ReadPagesByUserBetween :many
SELECT * FROM pages 
WHERE created BETWEEN sqlc.arg(start) AND sqlc.arg(end)
AND user_id = sqlc.arg(user_id)
ORDER BY created DESC;
