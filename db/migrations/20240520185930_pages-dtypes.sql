-- migrate:up
CREATE TABLE IF NOT EXISTS pages_new (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    url TEXT NOT NULL,
    title TEXT,
    description TEXT,
    image_url TEXT,
    preview_state TEXT DEFAULT 'unknown' NOT NULL,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

INSERT INTO pages_new (
    id,
    user_id,
    url,
    title,
    description,
    image_url,
    preview_state,
    created,
    updated
)
SELECT
    id,
    user_id,
    url,
    title,
    description,
    image_url,
    preview_state,
    created,
    updated
FROM pages;
DROP TABLE pages;
ALTER TABLE pages_new RENAME TO pages;


-- migrate:down
CREATE TABLE IF NOT EXISTS pages_new (
    id TEXT PRIMARY KEY NOT NULL,
    user_id TEXT NOT NULL,
    url TEXT NOT NULL,
    title TEXT,
    description TEXT,
    image_url TEXT,
    readability_status TEXT,
    readability_task_data TEXT,
    is_readable BOOL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL, preview_state TEXT DEFAULT 'unknown' NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users (id)
);
INSERT INTO pages_new (
    id,
    user_id,
    url,
    title,
    description,
    image_url,
    created,
    updated,
    preview_state
)
SELECT
    id,
    user_id,
    url,
    title,
    description,
    image_url,
    preview_state,
    created,
    updated
FROM pages;
DROP TABLE pages;
ALTER TABLE pages_new RENAME TO pages;
