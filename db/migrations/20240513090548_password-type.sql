-- migrate:up
CREATE TABLE IF NOT EXISTS "users_new" (
	id TEXT UNIQUE NOT NULL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL ,
	password BLOB NOT NULL,
	avatar TEXT,
	subscribed BOOL NOT NULL DEFAULT false,
	shortcut_token TEXT NOT NULL,
	has_readability BOOL NOT NULL DEFAULT false,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL
);
INSERT INTO users_new SELECT 
	id,
	username,
	email,
	CAST(password AS BLOB),
	avatar,
	subscribed,
	shortcut_token,
	has_readability,
	created,
	updated
	FROM users;
DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

-- migrate:down
CREATE TABLE "users_new" (
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
INSERT INTO users_new SELECT 
	id,
	username,
	email,
	CAST(password AS BINARY),
	avatar,
	subscribed,
	shortcut_token,
	has_readability,
	created,
	updated
	FROM users;
DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

