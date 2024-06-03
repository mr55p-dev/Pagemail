-- name: NewArticle :exec
INSERT INTO articles (page_id) VALUES (?);

-- name: UpdateArticleSuccess :exec
UPDATE articles
SET state = 'success',
	content = ?,
	updated = CURRENT_TIMESTAMP
WHERE page_id = ?;

-- name: UpdateArticleFailure :exec
UPDATE articles
SET state = 'failed',
	reason = ?,
	content = NULL,
	updated = CURRENT_TIMESTAMP
WHERE page_id = ?;
