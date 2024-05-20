package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	"github.com/mr55p-dev/pagemail/internal/tools"
)

func hashValue(val []byte) []byte {
	hasher := sha256.New()
	fmt.Fprint(hasher, val)

	out := make([]byte, 0)
	return hasher.Sum(out)
}

func NewResetToken() (token, hash []byte) {
	token = []byte(tools.GenerateNewId(64))
	hash = hashValue(token)
	return
}

func CheckResetToken(passed, stored []byte) bool {
	hashed := hashValue(passed)
	return subtle.ConstantTimeCompare(hashed, stored) == 1
}
