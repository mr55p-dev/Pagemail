package readability

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mr55p-dev/pagemail/internal/logging"
)

type service struct {
	baseUrl *url.URL
	client  http.Client
}

type optfn func(*service)

var logger = logging.NewLogger("readability")

func New(ctx context.Context, baseUrl *url.URL, opts ...optfn) (*service, error) {
	svc := &service{
		baseUrl: baseUrl,
		client:  *http.DefaultClient,
	}
	if err := svc.Ping(); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *service) Ping() error {
	pingUrl := s.baseUrl
	pingUrl = pingUrl.JoinPath("health")
	req, err := http.NewRequest(http.MethodGet, pingUrl.String(), nil)
	if err != nil {
		return err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Health check failed: %d %s", res.StatusCode, res.Status)
	}
	return nil
}
