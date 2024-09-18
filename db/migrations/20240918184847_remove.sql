-- migrate:up
DROP TABLE articles;
DROP TABLE readings;


-- migrate:down
CREATE TABLE articles (
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
CREATE TABLE readings (
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
