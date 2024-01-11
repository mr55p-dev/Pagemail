-- migrate:up
CREATE TABLE users (
	email STRING,
	password STRING
);

-- migrate:down
DROP TABLE users;
