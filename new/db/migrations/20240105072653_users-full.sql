-- migrate:up
DELETE FROM users;
DROP TABLE users;
CREATE TABLE users (
	id STRING UNIQUE NOT NULL PRIMARY KEY,
	username STRING,
	email STRING UNIQUE NOT NULL ,
	password STRING NOT NULL,
	name STRING,
	avatar STRING,
	subscribed BOOL DEFAULT false,
	shortcutToken STRING,
	has_readability BOOL DEFAULT false,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL
);

-- migrate:down
DROP TABLE users;
CREATE TABLE users (
	email STRING,
	password STRING
);
