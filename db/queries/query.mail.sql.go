// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.mail.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSchedule = `-- name: CreateSchedule :exec
INSERT INTO schedules (user_id, timezone, days, hour, minute) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type CreateScheduleParams struct {
	UserID   pgtype.UUID
	Timezone string
	Days     int32
	Hour     int32
	Minute   int32
}

func (q *Queries) CreateSchedule(ctx context.Context, arg CreateScheduleParams) error {
	_, err := q.db.Exec(ctx, createSchedule,
		arg.UserID,
		arg.Timezone,
		arg.Days,
		arg.Hour,
		arg.Minute,
	)
	return err
}

const readSchedules = `-- name: ReadSchedules :many
SELECT id, user_id, timezone, days, hour, minute, last_sent FROM schedules
`

func (q *Queries) ReadSchedules(ctx context.Context) ([]Schedule, error) {
	rows, err := q.db.Query(ctx, readSchedules)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Schedule
	for rows.Next() {
		var i Schedule
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Timezone,
			&i.Days,
			&i.Hour,
			&i.Minute,
			&i.LastSent,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateScheduleLastSent = `-- name: UpdateScheduleLastSent :exec
UPDATE schedules
SET last_sent = $1
WHERE user_id = $2
`

type UpdateScheduleLastSentParams struct {
	LastSent pgtype.Timestamp
	UserID   pgtype.UUID
}

func (q *Queries) UpdateScheduleLastSent(ctx context.Context, arg UpdateScheduleLastSentParams) error {
	_, err := q.db.Exec(ctx, updateScheduleLastSent, arg.LastSent, arg.UserID)
	return err
}
