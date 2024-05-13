-- name: CreateUser :exec
INSERT INTO users (
	id, username, email, password,
	avatar, subscribed, shortcut_token,
	has_readability, created, updated
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ReadUserById :one
SELECT * FROM users 
WHERE id = ? 
LIMIT 1;

-- name: ReadUserByShortcutToken :one
SELECT * FROM users 
WHERE shortcut_token = ?
LIMIT 1;

-- name: ReadUserByEmail :one
SELECT * FROM users 
WHERE email = ?
LIMIT 1;

-- name: ReadUsersWithMail :many
SELECT * FROM users 
WHERE subscribed = true;

-- name: ReadUserShortcutTokens :many
SELECT id, shortcut_token FROM users 
WHERE shortcut_token IS NOT NULL;

-- name: UpdateUser :exec
UPDATE users SET
	username = ?,
	password = ?,
	avatar = ?,
	subscribed = ?,
	shortcut_token = ?,
	has_readability = ?,
	updated = ?
WHERE id = :id;

-- name: UpdateUserPassword :exec
UPDATE users SET 
password = ? 
WHERE id = ?;

-- name: UpdateUserSubscription :exec
UPDATE users SET 
subscribed = ? 
WHERE id = ?;

-- name: UpdateUserShortcutToken :exec
UPDATE users SET 
shortcut_token = ? 
WHERE id = ?;

-- name: CreatePage :exec
INSERT INTO pages (id, user_id, url, created, updated)
VALUES (?, ?, ?, ?, ?);

-- name: ReadPageById :one
SELECT * FROM pages
WHERE id = ?
LIMIT 1;

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
AND user_id = ?;
-- ORDER BY created DESC;

-- name: UpsertPage :exec
INSERT OR REPLACE INTO pages (
	id, user_id, url, title, description,
	image_url, readability_status, readability_task_data,
	is_readable, created, updated
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: DeletePagesByUserId :execrows
DELETE FROM pages 
WHERE user_id = ?;
