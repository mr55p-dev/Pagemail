package mail

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/timer"
)

// Email address to send mail from
const (
	ADDR_DIGEST = "mail@pagemail.io"
	SUBJ_DIGEST = "Your daily digest"
	ADDR_RESET  = "support@pagemail.io"
	SUBJ_RESET  = "Reset your password"
)

var logger = logging.NewLogger("mail")

// Sender allows for different implementations of email clients using a send method
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// DB wraps the methods from database which are required for pulling users and their saved
// pages inside of an interval
type DB interface {
	ReadUsersWithMail(context.Context) ([]dbqueries.User, error)
	ReadPagesByUserBetween(context.Context, dbqueries.ReadPagesByUserBetweenParams) ([]dbqueries.Page, error)
}

// MailGo starts a goroutine on a timer to send emails to all subscribed users every 24 hours at
// 7 am
func MailGo(ctx context.Context, reader DB, sender Sender) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, time.Local)
	timer := timer.NewCronTimer(time.Hour*24, start)
	defer timer.Stop()

	for now := range timer.T {
		ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
		logger.InfoCtx(ctx, "Starting mail job", "time", now.Format(time.Stamp))
		err := DigestJob(ctx, reader, sender, now)
		cancel()
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to send digest", "error", err.Error())
		}
	}
}

// DigestJob collects all subscribed users, their pages between 24 hours ago and now, and then sends
// them
func DigestJob(ctx context.Context, reader DB, sender Sender, now time.Time) error {
	// Get users with mail enabled
	users, err := reader.ReadUsersWithMail(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get user list: %w", err)
	}

	// Dispatch jobs to the other goroutines
	var errCount int32
	jobs := make(chan dbqueries.User)
	errChan := make(chan error)
	wg := &sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range jobs {
				err := SendUserDigest(ctx, &user, reader, sender, now)
				if err != nil {
					errChan <- err
				}
			}
		}()
	}

	// Read out all the errors
	errList := make([]error, 0)
	go func() {
		for err := range errChan {
			atomic.AddInt32(&errCount, 1)
			errList = append(errList, err)
		}
	}()

	// Dispatch all the jobs
	for _, user := range users {
		jobs <- user
	}
	close(jobs)
	wg.Wait()

	if len(errList) > 0 {
		return errors.Join(errList...)
	}
	return nil
}

// SendUserDigest fetches the users pages, constructs an email and sends it via the sender interface
func SendUserDigest(ctx context.Context, user *dbqueries.User, db DB, sender Sender, now time.Time) error {
	logger := logger.With("user", user.Email)
	logger.DebugCtx(ctx, "Generating mail for user")

	// Read the users pages
	start := now.Add(-time.Hour * 24)
	end := now
	pages, err := db.ReadPagesByUserBetween(ctx, dbqueries.ReadPagesByUserBetweenParams{
		Start:  start,
		End:    end,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	// Brea
	if len(pages) == 0 {
		logger.DebugCtx(ctx, "No pages found")
		return nil
	}

	// Generate the mail and send it
	buf := bytes.Buffer{}
	err = render.MailDigest(now, user.Username, pages).Render(ctx, &buf)
	if err != nil {
		return err
	}
	msg := MakeMessage(
		user.Email,
		WithSubject(SUBJ_DIGEST),
		WithBody(&buf),
		WithTags(map[string]string{"puropse:": "digest"}),
	)
	err = sender.Send(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}
