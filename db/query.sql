-- name: GetUser :one
SELECT * FROM users 
WHERE id = ? 
LIMIT 1;

-- name: GetUserByShortcutToken :one
SELECT * FROM users
WHERE shortcut_token = ?
LIMIT 1;
