package tools

import (
	"crypto/rand"
	"math/big"
)

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenerateNewId(length int) string {
	// Generates a unique identifier of length chars from `alphabet`
	out := ""
	max := len(alphabet)
	for i := 0; i < length; i++ {
		char, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err != nil {
			panic(err)
		}
		out += string(alphabet[char.Int64()])
	}
	return out
}

func GenerateNewShortcutToken(id string) string {
	return "123"
}

