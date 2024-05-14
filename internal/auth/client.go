package auth

import (
	"context"
	"fmt"
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
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to get cookie: %w", err)
	}
	val, ok := a.store[cookie.Value]
	if !ok {
		return nil, fmt.Errorf("No such session: %w", err)
	}
	val.IsNew = false
	return val, nil
}

func (a *MemoryStore) New(r *http.Request, name string) (*sessions.Session, error) {
	tkn := tools.GenerateNewId(50)
	sess := &sessions.Session{
		ID:     tkn,
		Values: make(map[interface{}]interface{}),
		Options: &sessions.Options{
			MaxAge:   100,
			SameSite: http.SameSiteStrictMode,
		},
		IsNew: true,
	}
	a.store[name] = sess
	return sess, nil
}

func (a *MemoryStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	http.SetCookie(w, sessions.NewCookie(s.Name(), s.ID, s.Options))
	return nil
}

func GetLoginCookie(val string) *http.Cookie {
	return &http.Cookie{
		Name:     SessionKey,
		Value:    val,
		Path:     "/",
		MaxAge:   864000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}
