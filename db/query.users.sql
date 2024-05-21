-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    username,
    subscribed
) VALUES (?, ?, ?, ?)
RETURNING *;

-- name: ReadUserById :one
SELECT * FROM users 
WHERE id = ? 
LIMIT 1;

-- name: ReadUserByEmail :one
SELECT * FROM users 
WHERE email = ?
LIMIT 1;

-- name: ReadUsersWithMail :many
SELECT id, username, email FROM users 
WHERE subscribed = true;

-- name: UpdateUserSubscription :exec
UPDATE users SET 
subscribed = ? 
WHERE id = ?;
