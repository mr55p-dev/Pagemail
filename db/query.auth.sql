-- name: ReadUserByResetToken :one
SELECT * FROM auth LIMIT 1;

-- name: UpdateUserPassword :execrows
UPDATE auth
SET password_hash = ?
WHERE user_id = ?
    AND platform = 'pagemail'
RETURNING user_id;

-- -- name: UpdateUserShortcutToken :exec
UPDATE auth
SET shortcut_token = ?
WHERE user_id = ?
    AND platform = 'pagemail';

-- -- name: UpdateUserResetToken :exec
UPDATE auth
SET 
    reset_token = ?,
    reset_token_exp = ?
WHERE user_id = ?
    AND platform = 'pagemail';
