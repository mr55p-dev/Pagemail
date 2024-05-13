package tools

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"strings"
)

func GenerateNewId(length int) string {
	// Generates a unique identifier of length chars from `alphabet`
	out := new(strings.Builder)
	enc := base32.NewEncoder(base32.StdEncoding, out)
	defer enc.Close()

	n := length
	n /= 8
	if (length % 8) != 0 {
		n++
	}
	n *= 5

	// To generate a base32 encoding of 8 characters we must encode between 1 and 5 bytes
	written, err := io.CopyN(enc, rand.Reader, int64(n))
	_ = written
	if err != nil {
		panic(err)
	}
	enc.Close()
	return out.String()[:length]
}

func GenerateNewShortcutToken() string {
	return GenerateNewId(20)
}
