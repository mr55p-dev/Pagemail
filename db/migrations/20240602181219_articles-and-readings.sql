-- migrate:up
CREATE TABLE IF NOT EXISTS articles (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    page_id TEXT NOT NULL,
    state TEXT NOT NULL DEFAULT 'unknown',
    reason TEXT,
    content BLOB,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (page_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (page_id) REFERENCES pages (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS readings (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    article_id TEXT NOT NULL,
    job_id TEXT NOT NULL,
    state TEXT NOT NULL DEFAULT 'unknown',
    reason TEXT,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (article_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE CASCADE
);

ALTER TABLE pages DROP COLUMN reading_job_status;
ALTER TABLE pages DROP COLUMN reading_job_id;

-- migrate:down
ALTER TABLE pages ADD COLUMN reading_job_id TEXT;
ALTER TABLE pages ADD COLUMN reading_job_status TEXT NOT NULL DEFAULT 'unknown';
DROP TABLE readings;
DROP TABLE articles;
