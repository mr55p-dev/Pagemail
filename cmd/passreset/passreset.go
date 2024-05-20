package main

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mr55p-dev/gonk"
	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
)

var userId = flag.String("user-id", "", "User ID")

type Config struct {
	Url string `config:"db_url"`
}

func main() {
	ctx := context.Background()
	usage := func() {
		fmt.Fprintf(os.Stderr, "passreset can reset passwords in the database")
		flag.PrintDefaults()
	}
	flag.Usage = usage
	flag.Parse()

	cfg := new(Config)
	err := gonk.LoadConfig(cfg, gonk.EnvironmentLoader(""))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v", err)
		os.Exit(1)
	}

	if *userId == "" {
		fmt.Fprintf(os.Stderr, "user_id is required")
		usage()
		os.Exit(1)
	}

	password, err := parsePassword(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v", err)
		os.Exit(1)
	}
	passwordHash := auth.HashPassword(password)

	conn := db.MustConnect(ctx, cfg.Url)
	queries := dbqueries.New(conn)
	now := time.Now()
	_, hashedToken := auth.NewResetToken()
	queries.UpdateUserResetToken(ctx, dbqueries.UpdateUserResetTokenParams{
		ResetToken: hashedToken,
		ResetTokenExp: sql.NullTime{
			Time:  now.Add(time.Hour),
			Valid: true,
		},
		ID: *userId,
	})
	_, err = queries.UpdateUserPassword(ctx, dbqueries.UpdateUserPasswordParams{
		Password:      passwordHash,
		ResetToken:    hashedToken,
		ResetTokenExp: sql.NullTime{Valid: true, Time: now},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating password: %v", err)
		os.Exit(1)
	}
	fmt.Println("Password updated sucessfully")
	return
}

func parsePassword(rdr io.Reader) ([]byte, error) {
	password := new(bytes.Buffer)
	n, err := io.Copy(password, rdr)
	_ = n
	if err != nil {
		return nil, fmt.Errorf("Error reading password: %w", err)
	}
	passBytes := bytes.TrimSpace(password.Bytes())
	return passBytes, nil
}
