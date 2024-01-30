package driver

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// RestClient represents a REST client configuration.
type Client struct {
	Session            *Session
	BaseURL            string
	HTTPClient         *http.Client
	RequestReaderLimit int64
	// syncMutex  sync.Mutex // Mutex for ensuring thread safety
	DefaultExecuteStrategy CommandExecutor
}

type retryRoundTripper struct {
	next       http.RoundTripper
	maxRetries int
	delay      time.Duration
}

type logginRoundTripper struct {
	next http.RoundTripper
}

// RoundTrip
// retry decorator
func (rr retryRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// log.Println("pre-request in retry round tripper:", r.URL.Path)

	res, err := rr.next.RoundTrip(r)
	if err != nil {
		return res, err
	}

	// log.Println("post-response in retry round tripper:", res.StatusCode, res.Status)

	return res, nil
}

// RoundTrip
// logging decorator for each request to driver server
func (l logginRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Println(r.Method, r.URL.String())
	r.Header.Add("Accept", "json/application")
	res, err := l.next.RoundTrip(r)
	if err != nil {
		log.Println(err, res)
		return nil, err
	}

	return res, nil
}

// NewRestClient creates a new instance of the REST client with default settings.
func newClient(baseURL string, session *Session) *Client {
	return &Client{
		BaseURL:            baseURL,
		Session:            session,
		RequestReaderLimit: 200,
		HTTPClient: &http.Client{
			Transport: &retryRoundTripper{
				maxRetries: 3,
				delay:      time.Duration(200 * time.Millisecond),
				next: &logginRoundTripper{
					next: http.DefaultTransport,
				},
			},
		},
	}
}

func (cl Client) Execute(req *http.Request) (*http.Response, error) {
	return cl.HTTPClient.Do(req)
}

func (c Client) ExecuteCommandStrategy(cmd *Command, st ...CommandExecutor) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%s%s", c.BaseURL, c.Session.Route, c.Session.Id, cmd.Path)

	rr := io.LimitReader(ReusableReader(bytes.NewReader(cmd.Data)), 200)
	reqBody := io.NopCloser(rr)
	req, err := http.NewRequest(cmd.Method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if len(st) != 0 {
		for _, s := range st {
			return NewStrategy(s).Exec(req)
		}
	}

	return NewStrategy(c).Exec(req)

}
