package router

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gorilla/sessions"
)

const NumKeyBytes = 16

func loadCookieKey(router *Router, cfg *AppConfig) error {
	var input io.Reader
	if cfg.CookieKeyFile == "-" {
		input = rand.Reader
	} else {
		cookieDataFile, err := os.Open(cfg.CookieKeyFile)
		if err != nil {
			return fmt.Errorf("Failed to open cookie key file: %w", err)
		}
		defer cookieDataFile.Close()
		input = cookieDataFile
	}

	key := make([]byte, NumKeyBytes)
	n, err := input.Read(key)
	if err != nil {
		return err
	}
	if n < NumKeyBytes {
		return errors.New("Cookie key source has insufficient bytes")
	}
	router.Sessions = sessions.NewCookieStore(key)
	return nil
}
