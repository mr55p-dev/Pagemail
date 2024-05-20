-- name: ReadUserByResetToken :one
SELECT * FROM users
WHERE reset_token = ?
AND reset_token_exp > ?;

-- name: UpdateUserPassword :execrows
UPDATE users SET 
password = ? 
WHERE reset_token = ?
AND reset_token_exp > ?
RETURNING id;

-- name: UpdateUserShortcutToken :exec
UPDATE users SET 
shortcut_token = ? 
WHERE id = ?;

-- name: UpdateUserResetToken :exec
UPDATE users SET
reset_token = ?,
reset_token_exp = ?
WHERE id = ?;
