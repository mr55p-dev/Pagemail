-- migrate:up
CREATE TABLE IF NOT EXISTS "users_new" (
	id TEXT UNIQUE NOT NULL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL ,
	password BLOB NOT NULL,
	reset_token BLOB,
	reset_token_exp DATETIME,
	avatar TEXT,
	subscribed BOOL NOT NULL DEFAULT false,
	shortcut_token BLOB NOT NULL,
	has_readability BOOL NOT NULL DEFAULT false,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL
);
INSERT INTO users_new SELECT 
	id,
	username,
	email,
	password,
	reset_token,
	reset_token_exp,
	avatar,
	subscribed,
	CAST(shortcut_token AS BLOB),
	has_readability,
	created,
	updated
	FROM users;
DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

-- migrate:down
CREATE TABLE IF NOT EXISTS "users_new" (
	id TEXT UNIQUE NOT NULL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL ,
	password BLOB NOT NULL,
	reset_token BLOB,
	reset_token_exp DATETIME,
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
	password,
	reset_token,
	reset_token_exp,
	avatar,
	subscribed,
	CAST(shortcut_token AS TEXT),
	has_readability,
	created,
	updated
	FROM users;
DROP TABLE users;
ALTER TABLE users_new RENAME TO users;
