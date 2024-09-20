package main

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

// Mail sending timeout
const TIMEOUT = time.Second * 10

func openMailPool(username, password, host string, port, poolSize int) (*email.Pool, error) {
	mailAuth := smtp.PlainAuth("", username, password, host)
	connPool, err := email.NewPool(concatHostPort(host, port), poolSize, mailAuth)
	if err != nil {
		return nil, fmt.Errorf("Failed to open the connection pool: %w", err)
	}
	err = connPool.Send(&email.Email{
		To:      []string{"success@simulator.amazonses.com"},
		From:    formatAddress("Test Pagemail", "mail@pagemail.io"),
		Subject: "Test",
		Text:    []byte("Hello, world!"),
	}, TIMEOUT)
	if err != nil {
		return nil, fmt.Errorf("Failed to send the test email: %w", err)
	}
	return connPool, nil
}

func formatAddress(name string, address string) string {
	return fmt.Sprintf("%s <%s>", name, address)
}
