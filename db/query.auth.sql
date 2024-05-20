-- name: ReadUserByResetToken :one
SELECT users.* FROM auth
LEFT JOIN users
	ON users.id = auth.user_id
	AND auth.platform = 'pagemail'
WHERE auth.reset_token = ?
	AND reset_token_exp > ?;

-- name: UpdateUserPassword :execrows
UPDATE auth 
SET password_hash = ? 
FROM auth 
LEFT JOIN users
    ON auth.user_id = users.id
    AND auth.platform = 'pagemail'
WHERE auth.reset_token = ?
    AND auth.reset_token_exp > ?
RETURNING users.id;

-- name: UpdateUserShortcutToken :exec
UPDATE auth 
SET auth.shortcut_token = ? 
FROM auth
LEFT JOIN users
    ON users.id = auth.user_id
    AND auth.platform = 'pagemail'
WHERE users.id = ?;

-- name: UpdateUserResetToken :exec
UPDATE auth
SET 
    auth.reset_token = ?,
    auth.reset_token_exp = ?
FROM users
LEFT JOIN auth
	ON auth.user_id = users.id
	AND auth.platform = pagemail
WHERE users.id = ?;
