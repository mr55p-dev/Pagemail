package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mr55p-dev/pagemail/internal/mail"
)

var (
	recipient = flag.String("to", "mail@mail.com", "Address to send the test email to")
	subject   = flag.String("sub", "Test mail", "Subject of the test email")

	dry_run = flag.Bool("dry-run", false, "Weather to actually commit to sending the mail")
)

func main() {
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "mailmock [options] [message body]\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	body := flag.Arg(0)

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	client := mail.NewSesMailClient(ctx, cfg)
	bodyReader := strings.NewReader(body)
	if *dry_run {
		fmt.Println(body)
		fmt.Println("Dry-run enabled. Done.")
		return
	}
	err = client.Send(ctx, *recipient, bodyReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send message: \n%v", err)
		os.Exit(1)
	}
	fmt.Println("Sent message")

}
