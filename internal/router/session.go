package router

import (
	"errors"
	"io"

	"github.com/gorilla/sessions"
)

const NumKeyBytes = 16

func loadCookieKey(router *Router, input io.Reader) error {
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
