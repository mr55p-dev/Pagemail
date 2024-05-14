package auth

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

type MemoryStore struct {
	store map[string]*sessions.Session
}

func NewMemoryStore(ctx context.Context) *MemoryStore {
	return &MemoryStore{
		store: make(map[string]*sessions.Session),
	}
}

func (a *MemoryStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return nil, nil
}

func (a *MemoryStore) New(r *http.Request, name string) (*sessions.Session, error) {
	tkn := tools.GenerateNewId(50)
	sess := &sessions.Session{
		ID:      tkn,
		Values:  make(map[interface{}]interface{}),
		Options: &sessions.Options{},
		IsNew:   true,
	}
	a.store[name] = sess
	return sess, nil
}

func (a *MemoryStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return nil
}
