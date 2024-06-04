-- name: NewReading :one
INSERT INTO readings 
	(id, user_id, article_id, job_id, state) 
VALUES 
	(?, ?, ?, ?, ?)
RETURNING *;

-- name: GetReadingsByUser :many
SELECT * FROM readings
WHERE user_id = ?;

-- name: UpdateReading :exec
UPDATE readings
SET state = ?,
	reason = ?,
	updated = CURRENT_TIMESTAMP
