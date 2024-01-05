package db

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestUserRead(t *testing.T) {
	client := NewClient(zerolog.New(os.Stdout))
	user, err := client.ReadUserById("123")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(user, err)
}
