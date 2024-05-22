package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mr55p-dev/gonk"
	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
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
	err := gonk.LoadConfig(cfg, gonk.EnvLoader(""))
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
	_, err = queries.New(conn).UpdatePassword(ctx, queries.UpdatePasswordParams{
		PasswordHash: passwordHash,
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
