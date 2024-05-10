-- migrate:up
-- Remove the name column from users table
ALTER TABLE users
DROP COLUMN name;

-- Set the user_id field of pages to be NOT NULL
CREATE TABLE pages_new (
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
INSERT INTO pages_new SELECT * FROM pages;
DROP TABLE pages;
ALTER TABLE pages_new RENAME TO pages;

-- migrate:down
ALTER TABLE users
ADD COLUMN name TEXT;

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
