-- name: CreateLocalAuth :exec
INSERT INTO auth (
    user_id,
    platform,
	credential
) VALUES ($1, 'pagemail', crypt($2, gen_salt('bf')));

-- name: CreateShortcutAuth :exec
INSERT INTO auth (
    user_id,
    platform,
    credential
) VALUES ($1, 'shortcut', crypt($2, gen_salt('bf')));

-- name: CreateIdpAuth :exec
INSERT INTO auth (
    user_id,
    platform,
    credential
) VALUES ($1, $2, crypt($3, gen_salt('bf')));

-- name: ReadByUidPlatform :one
SELECT *
FROM auth
WHERE user_id = $1
AND platform = $2
LIMIT 1;

-- name: ReadAuthMethods :many
SELECT * FROM auth WHERE user_id = $1;

-- name: ReadUserByShortcut :one
SELECT users.*
FROM users
LEFT JOIN auth
    ON auth.user_id = users.id
    AND auth.platform = 'shortcut'
WHERE auth.credential = $1
LIMIT 1;

-- name: ReadByResetToken :one
SELECT user_id
FROM auth
WHERE reset_token = $1
    AND reset_expiry > $2
LIMIT 1;

-- name: UpdatePassword :execrows
UPDATE auth
SET credential = $1
WHERE user_id = $2
    AND platform = 'pagemail';

-- name: UpdateShortcutToken :exec
UPDATE auth
SET credential = $1
WHERE user_id = $2
    AND platform = 'shortcut';

-- name: UpdateResetToken :exec
UPDATE auth
SET 
    reset_token = $1,
    reset_expiry = $2
WHERE user_id = $3
    AND platform = 'pagemail';
