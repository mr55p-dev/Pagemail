package main

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/render"
)

// Mail sending timeout
const TIMEOUT = time.Second * 10

func openMailPool(username, password, host string, port, poolSize int) (*email.Pool, error) {
	mailAuth := smtp.PlainAuth("", username, password, host)
	connPool, err := email.NewPool(concatHostPort(host, port), poolSize, mailAuth)
	if err != nil {
		return nil, fmt.Errorf("Failed to open the connection pool: %w", err)
	}
	// err = connPool.Send(&email.Email{
	// 	To:      []string{"success@simulator.amazonses.com"},
	// 	From:    formatAddress("Test Pagemail", "mail@pagemail.io"),
	// 	Subject: "Test",
	// 	Text:    []byte("Hello, world!"),
	// }, TIMEOUT)
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to send the test email: %w", err)
	// }
	return connPool, nil
}

func formatAddress(name string, address string) string {
	return fmt.Sprintf("%s <%s>", name, address)
}

func getEmailContent(ctx context.Context, q *queries.Queries, from, until time.Time, userId string) ([]byte, error) {
	pages, err := q.ReadPagesByUserBetween(ctx, queries.ReadPagesByUserBetweenParams{
		Start:  from,
		End:    until,
		UserID: userId,
	})
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	for _, page := range pages {
		err = render.PageCard(page).Render(ctx, buf)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func MailJob2(ctx context.Context, pool *email.Pool, q *queries.Queries) error {
	// fetch the list of schedules
	schedules, err := q.ReadSchedules(ctx)
	if err != nil {
		return err
	}
	logger.InfoContext(ctx, "Collected schedules", "count", len(schedules))

	// check for schedules which need to be sent
	for _, schedule := range schedules {
		now := time.Now().UTC()
		loc, err := time.LoadLocation(schedule.Timezone)
		if err != nil {
			return err
		}
		var day int = now.Day()
		if schedule.Days != 0 && now.Weekday() != time.Weekday(schedule.Days)-1 {
			// it's not the right day of the week (some timezone fuckery definitely happens here)
			logger.InfoContext(ctx, "Not doing mail job - wrong day of week")
			continue
		}
		sendWindow := time.Date(
			now.Year(), now.Month(), day,
			int(schedule.Hour), int(schedule.Minute),
			0, 0, loc,
		)
		if now.Before(sendWindow) {
			// skip, we have not yet reached the cutoff window
			logger.InfoContext(ctx, "Not doing mail job - now before sendWindow", "sendWindow", sendWindow.Format(time.RFC3339))
			continue
		} else if schedule.LastSent.After(sendWindow) {
			// skip, we have sent an email corresponding to this schedule entry
			logger.InfoContext(
				ctx,
				"Not doing mail job - lastSent after sendWindow",
				"sendWindow", sendWindow.Format(time.RFC3339),
				"last sent", schedule.LastSent.Format(time.RFC3339),
			)
			continue
		}

		// load user
		user, err := q.ReadUserById(ctx, schedule.UserID)
		if err != nil {
			return fmt.Errorf("Could not find user for schedule")
		}

		// do send
		content, err := getEmailContent(ctx, q, schedule.LastSent, now, user.ID)
		if err != nil {
			return fmt.Errorf("Failed to get content: %w", err)
		}
		err = pool.Send(&email.Email{
			ReplyTo: []string{"mail@pagemail.io"},
			From:    formatAddress("Pagemail Daily Update", "mail@pagemail.io"),
			To:      []string{user.Email},
			Subject: "Your saved pages",
			HTML:    content,
		}, time.Second*3)
		if err != nil {
			return fmt.Errorf("Failed to send email: %w", err)
		}

		// update last sent
		err = q.UpdateScheduleLastSent(ctx, queries.UpdateScheduleLastSentParams{
			LastSent: now,
			UserID:   user.ID,
		})
		if err != nil {
			return fmt.Errorf("Failed to update last sent for mail")
		}

	}

	return nil
}

// func MailJob(ctx context.Context, pool *email.Pool, q *queries.Queries) {
// 	users, err := q.ReadMailUsersOverdue(ctx)
// 	if err != nil {
// 		PanicError("Could not read overdue uers", err)
// 	}
// 	logger.Info("Found overdue users", "count", len(users))
// 	grp := new(errgroup.Group)
// 	for _, user := range users {
// 		if user.ID != "DJTK42JZJ4" {
// 			continue
// 		}
// 		grp.Go(func() error {
// 			logger.Info("Starting mail job", "user", user.Email)
// 			until := time.Now().Round(time.Hour)
// 			from := time.Date(
// 				until.Year(),
// 				until.Month(),
// 				until.Day()-1,
// 				until.Hour(),
// 				until.Minute(),
// 				until.Second(),
// 				until.Nanosecond(),
// 				until.Location(),
// 			)
// 			content, err := getEmailContent(ctx, q, from, until, user.ID)
// 			if err != nil {
// 				return fmt.Errorf("Failed to get content: %w", err)
// 			}
// 			err = pool.Send(&email.Email{
// 				ReplyTo: []string{"mail@pagemail.io"},
// 				From:    formatAddress("Pagemail Daily Update", "mail@pagemail.io"),
// 				To:      []string{user.Email},
// 				Subject: "Your saved pages",
// 				HTML:    content,
// 			}, time.Second*3)
// 			if err != nil {
// 				return fmt.Errorf("Failed to send email: %w", err)
// 			}
//
// 			err = q.UpdateMailNextTS(ctx, queries.UpdateMailNextTSParams{
// 				NextMailTs: sql.NullTime{Time: until.Add(time.Hour * 24), Valid: true},
// 				ID:         user.ID,
// 			})
// 			if err != nil {
// 				return fmt.Errorf("Failed to update next timestamp: %w", err)
// 			}
//
// 			return nil
// 		})
// 	}
// 	err = grp.Wait()
// 	if err != nil {
// 		PanicError("something threw", err)
// 	}
// }
