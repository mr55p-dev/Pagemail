-- migrate:up
ALTER TABLE users ADD COLUMN reset_token BLOB;
ALTER TABLE users ADD COLUMN reset_token_exp DATETIME;


-- migrate:down
ALTER TABLE users DROP COLUMN reset_token;
ALTER TABLE users DROP COLUMN reset_token_exp;
