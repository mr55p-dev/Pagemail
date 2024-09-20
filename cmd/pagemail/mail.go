package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/smtp"
	"sync"
	"time"

	"github.com/jordan-wright/email"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/render"
)

type sch struct {
	sync.RWMutex
	items []queries.Schedule
}

type Mailer struct {
	schedules *sch
	timeout   time.Duration
	pool      *email.Pool
	conn      *sql.DB
}

// New will create a new mailer and populate it's schedules list. The self populatio ends when the context is cancelled
func New(ctx context.Context, db *sql.DB, pool *email.Pool, timeout time.Duration) *Mailer {

	m := &Mailer{
		schedules: &sch{
			RWMutex: sync.RWMutex{},
			items:   []queries.Schedule{},
		},
		timeout: timeout,
		pool:    pool,
		conn:    db,
	}
	m.updateSchedules(ctx)

	go func() {
		for {
			timer := time.NewTimer(time.Minute * 10)
			select {
			case <-timer.C:
				m.updateSchedules(ctx)
			case <-ctx.Done():
				timer.Stop()
				return
			}
		}
	}()

	return m
}

// updateSchedules updates the in memory schedule representation
func (m *Mailer) updateSchedules(ctx context.Context) error {
	m.schedules.Lock()
	defer m.schedules.Unlock()

	schedules, err := m.Queries().ReadSchedules(ctx)
	if err != nil {
		return err
	}

	m.schedules.items = schedules
	return nil
}

// StartMailJob will run the scheduled mailer at an interval specified, stopping once the context is cancelled
func (m *Mailer) StartMailJob(ctx context.Context, interval time.Duration) {
	go func() {
		for {
			timer := time.NewTimer(interval)
			select {
			case now := <-timer.C:
				childCtx, _ := context.WithTimeout(ctx, interval)
				m.RunScheduledSend(childCtx, now)
			case <-ctx.Done():
				timer.Stop()
				return
			}

		}
	}()
}

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

func (m *Mailer) RunScheduledSend(ctx context.Context, now time.Time) error {
	m.schedules.RLock()
	defer m.schedules.RUnlock()

	// check for schedules which need to be sent
	for _, schedule := range m.schedules.items {
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
	}

	return nil
}

func (m *Mailer) Queries() *queries.Queries {
	return queries.New(m.conn)
}

// sendForSchedule will send the email required by a schedule
func (m *Mailer) sendForSchedule(ctx context.Context, schedule queries.Schedule, now time.Time) error {
	// load user
	user, err := m.Queries().ReadUserById(ctx, schedule.UserID)
	if err != nil {
		return fmt.Errorf("Could not find user for schedule")
	}

	// load pages
	pages, err := m.Queries().ReadPagesByUserBetween(ctx, queries.ReadPagesByUserBetweenParams{
		Start:  schedule.LastSent,
		End:    now,
		UserID: schedule.UserID,
	})
	if err != nil {
		return fmt.Errorf("Failed to load pages for user %s: %w", schedule.UserID, err)
	}

	// if len(pages) == 0 {
	// 	return nil
	// }

	// send the email
	content, err := getEmailContent(ctx, &user, pages)
	if err != nil {
		return fmt.Errorf("Failed to get content: %w", err)
	}
	err = m.pool.Send(content, m.timeout)
	if err != nil {
		return fmt.Errorf("Failed to send email: %w", err)
	}

	// update last sent
	err = m.Queries().UpdateScheduleLastSent(ctx, queries.UpdateScheduleLastSentParams{
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
