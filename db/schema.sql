CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE users (
	id STRING UNIQUE NOT NULL PRIMARY KEY,
	username STRING,
	email STRING UNIQUE NOT NULL ,
	password BINARY NOT NULL,
	name STRING,
	avatar STRING,
	subscribed BOOL DEFAULT false,
	shortcut_token STRING,
	has_readability BOOL DEFAULT false,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL
);
CREATE TABLE pages (
    id STRING PRIMARY KEY NOT NULL,
    user_id STRING,
    url STRING NOT NULL,
    title STRING,
    description STRING,
    image_url STRING,
    readability_status STRING,
    readability_task_data STRING,
    is_readable BOOL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,

	FOREIGN KEY(user_id) REFERENCES users(id)
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20240104200335'),
  ('20240105072653'),
  ('20240105122600');
