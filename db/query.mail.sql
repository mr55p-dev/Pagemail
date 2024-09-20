-- name: CreateSchedule :exec
INSERT INTO schedules (user_id, timezone, days, hour, minute) 
VALUES (?, ?, ?, ?, ?)
RETURNING id;

-- name: ReadSchedules :many
SELECT * FROM schedules;

-- name: UpdateScheduleLastSent :exec
UPDATE schedules
SET last_sent = ?
WHERE user_id = ?;
