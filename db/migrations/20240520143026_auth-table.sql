-- migrate:up
CREATE TABLE IF NOT EXISTS users_new (
    id TEXT PRIMARY KEY NOT NULL,
    email TEXT NOT NULL,
    username TEXT NOT NULL,
	subscribed BOOL NOT NULL DEFAULT FALSE,
    has_readability BOOL NOT NULL DEFAULT FALSE,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS auth (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT,
    platform TEXT NOT NULL,
    password_hash BLOB,
    password_reset_token BLOB,
    password_reset_expiry DATETIME,
    credential BLOB,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users_new (id) ON DELETE CASCADE,

    UNIQUE (user_id, platform),
    UNIQUE (platform, credential)
);

INSERT INTO auth (user_id, platform, password_hash, created, updated)
SELECT
    id,
    'pagemail',
    password,
    created,
    updated
FROM users;

INSERT INTO auth (user_id, platform, credential)
SELECT 
	id, 
	'shortcut',
	shortcut_token
FROM users;



INSERT INTO users_new (id, email, username, subscribed, created, updated)
SELECT
    id,
    email,
    username,
	subscribed,
    created,
    updated
FROM users;
DROP TABLE users;
ALTER TABLE users_new RENAME TO users;


-- migrate:down
CREATE TABLE IF NOT EXISTS users_new (
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

INSERT INTO users_new (id, username, email, password, subscribed, shortcut_token, has_readability, created, updated)
	SELECT 
		user.id,
		user.username,
		user.email,
		auth.password_hash,
		user.subscribed,
		shortcut.credential,
		FALSE,
		user.created,
		user.updated
	FROM users as user
	LEFT JOIN auth
	ON user.id = auth.user_id
	AND auth.platform = 'pagemail'
	LEFT JOIN auth as shortcut
	ON shortcut.user_id = user.id
	AND shortcut.platform = 'shortcut';


DROP TABLE users;
ALTER TABLE users_new RENAME TO users;
