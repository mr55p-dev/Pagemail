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

-- name: ReadUserWithCredential :one
SELECT users.* 
FROM users
LEFT JOIN auth 
ON users.id = auth.user_id
WHERE users.email = $1
AND auth.platform = $2
AND auth.credential = crypt($3, auth.credential);
