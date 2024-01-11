-- migrate:up
CREATE TABLE pages (
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

-- migrate:down
DROP TABlE pages;
