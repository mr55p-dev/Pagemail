-- migrate:up
CREATE TABLE users_new (
	id TEXT UNIQUE NOT NULL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL ,
	password BINARY NOT NULL,
	name TEXT,
	avatar TEXT,
	subscribed BOOL NOT NULL DEFAULT false,
	shortcut_token TEXT NOT NULL,
	has_readability BOOL NOT NULL DEFAULT false,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL
);
INSERT INTO users_new SELECT * FROM users;
DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

CREATE TABLE pages_new (
	id TEXT PRIMARY KEY NOT NULL,
	user_id TEXT,
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
INSERT INTO pages_new SELECT * FROM pages;
DROP TABLE pages;
ALTER TABLE pages_new RENAME TO pages;

-- migrate:down
CREATE TABLE users_new (
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
INSERT INTO users_new SELECT * FROM users;
DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

CREATE TABLE pages_new (
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
INSERT INTO pages_new SELECT * FROM pages;
DROP TABLE pages;
ALTER TABLE pages_new RENAME TO pages;
