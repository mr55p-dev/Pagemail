package router

import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/gorilla/sessions"
)

func loadCookieKey(router *Router, input io.Reader) error {
	key, err := io.ReadAll(input)
	if err != nil {
		return err
	}

	// TEMP
	// try to base64 decode the key
	// if it fails, it's not base64 encoded
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(key))
	_, err = io.ReadAll(dec)
	if err != nil {
		panic(err)
	}

	router.Sessions = sessions.NewCookieStore(key)
	return nil
}
