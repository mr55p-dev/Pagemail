-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    username
) VALUES ($1, $2, $3)
RETURNING *;

-- name: ReadUserById :one
SELECT * FROM users 
WHERE id = $1 
LIMIT 1;

-- name: ReadUserByEmail :one
SELECT * FROM users 
WHERE email = $1
LIMIT 1;
