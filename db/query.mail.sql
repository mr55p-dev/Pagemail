-- name: CreateSchedule :exec
INSERT INTO schedules (user_id, timezone, days, hour, minute) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: ReadSchedules :many
SELECT * FROM schedules;

-- name: UpdateScheduleLastSent :exec
UPDATE schedules
SET last_sent = $1
WHERE user_id = $2;
