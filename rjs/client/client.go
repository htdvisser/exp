package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"htdvisser.dev/exp/rjs"
)

// New returns a new RJS Client on top of the given HTTP client.
// It prepends the given base URL to API URIs, and sets the given Authorization header if not empty.
func New(client *http.Client, baseURL, authorization string) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	base.RawQuery = ""
	base.Path = strings.TrimSuffix(base.Path, "/")
	return &Client{
		client:        client,
		baseURL:       base.String(),
		authorization: authorization,
	}, nil
}

// Client is an HTTP client that calls the RJS API.
type Client struct {
	client        *http.Client
	baseURL       string
	authorization string
}

func (c *Client) url(uri string) string {
	return c.baseURL + uri
}

// Call implements the rjs.Client interface.
func (c *Client) Call(ctx context.Context, uri string, request interface{}, response interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return fmt.Errorf("failed to encode request body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url(uri), &buf)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf(
		"%s go/%s",
		filepath.Base(os.Args[0]),
		strings.TrimPrefix(runtime.Version(), "go"),
	))
	if c.authorization != "" {
		req.Header.Set("Authorization", c.authorization)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return &rjs.Error{
			StatusCode: res.StatusCode,
			Body:       body,
		}
	}
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}
	return nil
}
