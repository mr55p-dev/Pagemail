CREATE TABLE IF NOT EXISTS schema_migrations (version varchar(128) PRIMARY KEY);
CREATE TABLE IF NOT EXISTS users (
    id text PRIMARY KEY NOT NULL,
    email text NOT NULL,
    username text NOT NULL,
    subscribed bool NOT NULL DEFAULT FALSE,
    has_readability bool NOT NULL DEFAULT FALSE,
    created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE auth (
    id integer PRIMARY KEY AUTOINCREMENT,
    user_id text NOT NULL,
    platform text NOT NULL,
    password_hash blob,
    password_reset_token blob,
    password_reset_expiry datetime,
    credential blob,
    created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,

    UNIQUE (user_id, platform),
    UNIQUE (platform, credential)
);
CREATE TABLE IF NOT EXISTS pages (
    id text PRIMARY KEY,
    user_id text NOT NULL,
    url text NOT NULL,
    title text,
    description text,
    image_url text,
    preview_state text DEFAULT 'unknown' NOT NULL,
    created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    readable bool NOT NULL DEFAULT FALSE,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
CREATE TABLE articles (
    id integer PRIMARY KEY AUTOINCREMENT,
    page_id text NOT NULL,
    state text NOT NULL DEFAULT 'unknown',
    reason text,
    html blob NOT NULL,
    content blob NOT NULL,
    created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (page_id),
    FOREIGN KEY (page_id) REFERENCES pages (id) ON DELETE CASCADE
);
CREATE TABLE readings (
    id integer PRIMARY KEY AUTOINCREMENT,
    article_id text NOT NULL,
    job_id text,
    state text NOT NULL DEFAULT 'unknown',
    reason text,
    created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (article_id),
    FOREIGN KEY (id) REFERENCES articles (id) ON DELETE CASCADE
);
-- Dbmate schema migrations
INSERT INTO schema_migrations (version) VALUES
('20240104200335'),
('20240105072653'),
('20240105122600'),
('20240510132816'),
('20240510134137'),
('20240513090548'),
('20240515151753'),
('20240520062158'),
('20240520065956'),
('20240520143026'),
('20240520185930'),
('20240522124057'),
('20240602181219');
