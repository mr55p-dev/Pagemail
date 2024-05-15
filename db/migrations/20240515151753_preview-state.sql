-- migrate:up
ALTER TABLE pages ADD COLUMN preview_state TEXT DEFAULT 'unknown' NOT NULL;

-- migrate:down
ALTER TABLE pages DROP COLUMN preview_state;

