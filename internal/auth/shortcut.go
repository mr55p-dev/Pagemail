package auth

import "github.com/mr55p-dev/pagemail/internal/tools"

func NewShortcutToken() (token, hash []byte) {
	token = []byte(tools.GenerateNewId(32))
	hash = HashValue(token)
	return
}
