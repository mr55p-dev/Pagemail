package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mr55p-dev/pagemail/internal/auth"
)

var user_id = flag.String("user-id", "", "User ID")

func main() {
	usage := func() {
		fmt.Fprintf(os.Stderr, "passreset can reset passwords in the database")
		flag.PrintDefaults()
	}
	flag.Usage = usage
	flag.Parse()

	if *user_id == "" {
		fmt.Fprintf(os.Stderr, "user_id is required")
		usage()
		os.Exit(1)
	}

	password := new(bytes.Buffer)
	io.Copy(os.Stdin, password)
	passwordHash := auth.HashPassword(password.Bytes())
	fmt.Fprint(os.Stdout, string(passwordHash))
}
