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

-- name: CreateIdpAuth :exec
INSERT INTO auth (
    user_id,
    platform,
    credential
) VALUES (?, ?, ?);

-- name: ReadByUidPlatform :one
SELECT *
FROM auth
WHERE user_id = ?
AND platform = ?
LIMIT 1;

-- name: ReadAuthMethods :many
SELECT * FROM auth WHERE user_id = ?;

-- name: ReadUserByShortcut :one
SELECT users.*
FROM users
LEFT JOIN auth
    ON auth.user_id = users.id
    AND auth.platform = 'shortcut'
WHERE auth.credential = ?
LIMIT 1;

-- name: ReadByResetToken :one
SELECT user_id
FROM auth
WHERE password_reset_token = ?
    AND password_reset_expiry > ?
LIMIT 1;

-- name: UpdatePassword :execrows
UPDATE auth
SET password_hash = ?
WHERE user_id = ?
    AND platform = 'pagemail';

-- name: UpdateShortcutToken :exec
UPDATE auth
SET credential = ?
WHERE user_id = ?
    AND platform = 'shortcut';

-- name: UpdateResetToken :exec
UPDATE auth
SET 
    password_reset_token = ?,
    password_reset_expiry = ?
WHERE user_id = ?
    AND platform = 'pagemail';
