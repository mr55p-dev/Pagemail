package router

import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func loadCookieKey(router *Router, input io.Reader) error {
	key, err := io.ReadAll(input)
	if err != nil {
		return err
	}

	// TEMP
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(key))
	_, err = io.ReadAll(dec)
	if err != nil {
		panic(err)
	}

	_ = key
	router.Sessions = sessions.NewCookieStore(securecookie.GenerateRandomKey(20))
	return nil
}
