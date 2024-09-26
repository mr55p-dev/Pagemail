-- migrate:up
ALTER TABLE users DROP COLUMN next_mail_ts;
CREATE TABLE schedules (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    timezone TEXT NOT NULL DEFAULT 'Europe/London',
    days INTEGER NOT NULL DEFAULT 0,
    hour INTEGER NOT NULL,
    minute INTEGER NOT NULL DEFAULT 0,
    last_sent DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- CREATE TABLE emails (
-- 	id INTEGER PRIMARY KEY AUTOINCREMENT,
-- 	user_id TEXT NOT NULL,
-- 	schedule_id TEXT NOT NULL
-- 	status STRING NOT NULL DEFAULT 'pending',
-- 	
-- 	FOREIGN KEY (user_id) REFERENCES users (id),
-- 	FOREIGN KEY (schedule_id) REFERENCES schedules (id)
-- );

-- migrate:down
DROP TABLE schedules;
ALTER TABLE users ADD COLUMN next_mail_ts DATETIME;
