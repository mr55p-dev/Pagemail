-- migrate:up
ALTER TABLE users ADD COLUMN next_mail_ts DATETIME;

-- migrate:down
ALTER TABLE users DROP COLUMN next_mail_ts;
