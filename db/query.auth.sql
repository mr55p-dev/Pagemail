-- name: CreateLocalAuth :exec
INSERT INTO auth (
    user_id,
    platform,
    password_hash
) VALUES (?, 'pagemail', ?);

-- name: CreateShortcutAuth :exec
INSERT INTO auth (
    user_id,
    platform,
    credential
) VALUES (?, 'shortcut', ?);

-- name: ReadUserByResetToken :one
SELECT user_id
FROM auth
WHERE password_reset_token = ?
    AND password_reset_expiry > ?
LIMIT 1;

-- name: UpdateUserPassword :execrows
UPDATE auth
SET password_hash = ?
WHERE user_id = ?
    AND platform = 'pagemail'
RETURNING user_id;

-- name: UpdateUserShortcutToken :exec
UPDATE auth
SET credential = ?
WHERE user_id = ?
    AND platform = 'shortcut';

-- -- name: UpdateUserResetToken :exec
UPDATE auth
SET 
    reset_token = ?,
    reset_token_exp = ?
WHERE user_id = ?
    AND platform = 'pagemail';
