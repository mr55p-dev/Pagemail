package readability

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mr55p-dev/pagemail/internal/logging"
)

type Client struct {
	baseUrl *url.URL
	client  http.Client
}

type optfn func(*Client)

var logger = logging.NewLogger("readability")

func New(ctx context.Context, baseUrl *url.URL, opts ...optfn) *Client {
	svc := &Client{
		baseUrl: baseUrl,
		client:  *http.DefaultClient,
	}
	return svc
}

func (s *Client) Ping() error {
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

func (s *Client) Check(ctx context.Context, url string, content io.Reader) (bool, error) {
	// construct url
	checkUrl := s.baseUrl
	checkUrl = checkUrl.JoinPath("check")
	q := checkUrl.Query()
	q.Set("url", url)
	checkUrl.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, checkUrl.String(), content)
	if err != nil {
		return false, err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Check failed: %d %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("Failed to read response body: %w", err)
	}
	logger.InfoCtx(ctx, "Readability check response", "res", resBytes)

	resJson := &struct {
		IsReadable bool `json:"is_readable"`
	}{}
	err = json.Unmarshal(resBytes, resJson)
	if err != nil {
		return false, fmt.Errorf("Failed to parse response: %w", err)
	}

	return resJson.IsReadable, nil
}
