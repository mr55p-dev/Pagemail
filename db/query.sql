-- name: CreateUser :exec
INSERT INTO users (
	id, username, email, password, name,
	avatar, subscribed, shortcut_token,
	has_readability, created, updated
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);


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
	name = ?,
	avatar = ?,
	subscribed = ?,
	shortcut_token = ?,
	has_readability = ?,
	updated = ?
WHERE id = :id
