package mail

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/render"
)

type Mailer struct {
	Timetout time.Duration
	Pool     *email.Pool
	DB       *sql.DB
}

func NewPool(username, password, host string, port, poolSize int) (*email.Pool, error) {
	mailAuth := smtp.PlainAuth("", username, password, host)
	connPool, err := email.NewPool(fmt.Sprintf("%s:%d", host, port), poolSize, mailAuth)
	if err != nil {
		return nil, fmt.Errorf("Failed to open the connection pool: %w", err)
	}
	err = connPool.Send(&email.Email{
		To:      []string{"success@simulator.amazonses.com"},
		From:    formatAddress("Test Pagemail", "mail@pagemail.io"),
		Subject: "Test",
		Text:    []byte("Hello, world!"),
	}, time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("Failed to send the test email: %w", err)
	}
	return connPool, nil
}

func (m *Mailer) RunScheduledSend(ctx context.Context, now time.Time) (int, error) {
	schedules, err := m.queries().ReadSchedules(ctx)
	count := 0
	if err != nil {
		return count, fmt.Errorf("Failed to get schedules: %w", err)
	}

	// check for schedules which need to be sent
	errList := []error{}
	for _, schedule := range schedules {
		loc, err := time.LoadLocation(schedule.Timezone)
		if err != nil {
			return count, err
		}
		var day int = now.Day()
		if schedule.Days != 0 && now.Weekday() != time.Weekday(schedule.Days)-1 {
			// it's not the right day of the week (some timezone fuckery definitely happens here)
			continue
		}
		sendWindow := time.Date(
			now.Year(), now.Month(), day,
			int(schedule.Hour), int(schedule.Minute),
			0, 0, loc,
		)
		if now.Before(sendWindow) {
			// skip, we have not yet reached the cutoff window
			continue
		} else if schedule.LastSent.After(sendWindow) {
			// skip, we have sent an email corresponding to this schedule entry
			continue
		}

		if err := m.sendForSchedule(ctx, schedule, now); err != nil {
			errList = append(errList, err)
		} else {
			count++
		}
	}
	return count, errors.Join(errList...)
}

func (m *Mailer) queries() *queries.Queries {
	return queries.New(m.conn)
}

// sendForSchedule will send the email required by a schedule
func (m *Mailer) sendForSchedule(ctx context.Context, schedule queries.Schedule, now time.Time) error {
	// load user
	user, err := m.queries().ReadUserById(ctx, schedule.UserID)
	if err != nil {
		return fmt.Errorf("Could not find user for schedule")
	}

	// load pages
	pages, err := m.queries().ReadPagesByUserBetween(ctx, queries.ReadPagesByUserBetweenParams{
		Created:   schedule.LastSent,
		Created_2: now,
		UserID:    schedule.UserID,
	})
	if err != nil {
		return fmt.Errorf("Failed to load pages for user %s: %w", schedule.UserID, err)
	}

	if len(pages) == 0 {
		return nil
	}

	// send the email
	content, err := getEmailContent(ctx, &user, pages)
	if err != nil {
		return fmt.Errorf("Failed to get content: %w", err)
	}

	err = m.Pool.Send(content, m.Timetout)
	if err != nil {
		return fmt.Errorf("Failed to send email: %w", err)
	}

	// update last sent
	err = m.queries().UpdateScheduleLastSent(ctx, queries.UpdateScheduleLastSentParams{
		LastSent: now,
		UserID:   user.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to update last sent for mail")
	}

	return nil
}

// getEmailContent will compile an email to send
func getEmailContent(ctx context.Context, user *queries.User, pages []queries.Page) (*email.Email, error) {
	buf := new(bytes.Buffer)
	for _, page := range pages {
		err := render.PageCard(page).Render(ctx, buf)
		if err != nil {
			return nil, err
		}
	}
	return &email.Email{
		ReplyTo: []string{"mail@pagemail.io"},
		From:    formatAddress("Pagemail Daily Update", "mail@pagemail.io"),
		To:      []string{user.Email},
		Subject: "Your saved pages",
		HTML:    buf.Bytes(),
	}, nil

}

// formatAddress will give a display name to an email address
func formatAddress(name string, address string) string {
	return fmt.Sprintf("%s <%s>", name, address)
}
