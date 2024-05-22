-- migrate:up
ALTER TABLE pages ADD COLUMN readable BOOL NOT NULL DEFAULT false;
ALTER TABLE pages ADD COLUMN reading_job_status STRING NOT NULL DEFAULT 'unknown';
ALTER TABLE pages ADD COLUMN reading_job_id STRING;


-- migrate:down
ALTER TABLE pages DROP COLUMN readable;
ALTER TABLE pages DROP COLUMN reading_job_status;
ALTER TABLE pages DROP COLUMN reading_job_id;
