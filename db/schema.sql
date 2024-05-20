CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE IF NOT EXISTS "pages" (
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

	FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS "users" (
    id TEXT PRIMARY KEY NOT NULL,
    email TEXT NOT NULL,
    username TEXT NOT NULL,
	subscribed BOOL NOT NULL DEFAULT FALSE,
    has_readability BOOL NOT NULL DEFAULT FALSE,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE auth (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT,
    platform TEXT NOT NULL,
    password_hash BLOB,
    password_reset_token BLOB,
    password_reset_expiry DATETIME,
    credential BLOB,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES "users" (id) ON DELETE CASCADE,

    UNIQUE (user_id, platform),
    UNIQUE (platform, credential)
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20240104200335'),
  ('20240105072653'),
  ('20240105122600'),
  ('20240510132816'),
  ('20240510134137'),
  ('20240513090548'),
  ('20240515151753'),
  ('20240520062158'),
  ('20240520065956'),
  ('20240520143026');
