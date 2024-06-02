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

type JobError struct {
	Message string `json:"msg"`
	Detail  string `json:"detail"`
}

type Status struct {
	JobId  string     `json:"jobId"`
	Status string     `json:"status"`
	Reason string     `json:"reason"`
	Errors []JobError `json:"errors"`
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
	req.Header.Add("Content-Type", "text/html")
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

func (s *Client) Extract(ctx context.Context, url string, content io.Reader) ([]byte, error) {
	dest := s.baseUrl.JoinPath("/extract")
	q := dest.Query()
	q.Add("url", url)
	dest.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, dest.String(), content)
	if err != nil {
		return nil, fmt.Errorf("Failed to construct request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Non 200 response code: %s", res.Status)
	}

	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}
	return buf, nil
}

func (s *Client) Synthesize(ctx context.Context, textContent io.Reader) (*Status, error) {
	reqUrl := s.baseUrl.JoinPath("/synthesize")
	req, err := http.NewRequest(http.MethodPost, reqUrl.String(), textContent)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "text/html")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get response: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Non-200 status code: %s", res.Status)
	}

	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	rval := new(Status)
	err = json.Unmarshal(content, rval)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %w", err)
	}

	return rval, nil
}

func (s *Client) Status(ctx context.Context, id string) (*Status, error) {
	reqUrl := s.baseUrl.JoinPath("/status")
	q := reqUrl.Query()
	q.Add("jobId", id)
	reqUrl.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, reqUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get response: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Non-200 status code: %s", res.Status)
	}

	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	rval := new(Status)
	err = json.Unmarshal(content, rval)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %w", err)
	}

	return rval, nil
}
