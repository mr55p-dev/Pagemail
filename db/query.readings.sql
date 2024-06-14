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
	updated = CURRENT_TIMESTAMP;

-- name: GetAllReadingInfo :many
SELECT sqlc.embed(pages), sqlc.embed(articles), sqlc.embed(readings)
FROM pages
INNER JOIN articles
ON page.id = article.page_id
LEFT JOIN readings
ON article.id = reading.article_id
WHERE pages.user_id = ?;
