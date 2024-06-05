-- name: NewArticle :exec
INSERT INTO articles (id, user_id, page_id) VALUES (?, ?, ?);

-- name: GetArticle :one
SELECT * FROM articles WHERE id = ?;

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

-- name: GetPagesAndArticles :many
SELECT sqlc.embed(pages), sqlc.embed(articles)
FROM pages
INNER JOIN articles
ON pages.id = articles.page_id
WHERE pages.user_id = ?;
