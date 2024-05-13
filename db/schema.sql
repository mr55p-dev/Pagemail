CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE IF NOT EXISTS "users" (
	id TEXT UNIQUE NOT NULL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL ,
	password BINARY NOT NULL,
	avatar TEXT,
	subscribed BOOL NOT NULL DEFAULT false,
	shortcut_token TEXT NOT NULL,
	has_readability BOOL NOT NULL DEFAULT false,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL
);
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
	updated DATETIME NOT NULL,

	FOREIGN KEY(user_id) REFERENCES users(id)
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20240104200335'),
  ('20240105072653'),
  ('20240105122600'),
  ('20240510132816'),
  ('20240510134137');
