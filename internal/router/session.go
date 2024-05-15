package router

import (
	"io"

	"github.com/gorilla/sessions"
)

func loadCookieKey(router *Router, input io.Reader) error {
	key, err := io.ReadAll(input)
	if err != nil {
		return err
	}

	router.Sessions = sessions.NewCookieStore(key)
	return nil
}
