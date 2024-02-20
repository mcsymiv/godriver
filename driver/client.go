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

// retryRoundTripper
// http.RoundTrip client middleware
// TODO: add retry logic
// most of request retries will be handled in strategies
// serves as a general retry between client and webdriver
type retryRoundTripper struct {
	next       http.RoundTripper
	maxRetries int
	delay      time.Duration
}

// RoundTrip
// middleware for retries
// TODO: add global retry logic
func (rr retryRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := rr.next.RoundTrip(r)
	if err != nil {
		return res, err
	}

	return res, nil
}

// loggingRoundTripper
type logginRoundTripper struct {
	next http.RoundTripper
}

// RountTrip
// middleware logger for Client
func (l logginRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Printf("req: %s %s", r.Method, r.URL.Path)
	res, err := l.next.RoundTrip(r)
	if err != nil {
		return nil, fmt.Errorf("error on %v request: %v", r, err)
	}

	log.Printf("res status: %v", res.StatusCode)
	return res, nil
}

// newClient
// NewRestClient creates a new instance of the REST client with default settings.
func newClient(baseURL string, session *Session) *Client {
	return &Client{
		BaseURL:            baseURL,
		Session:            session,
		RequestReaderLimit: 2048,
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

// Execute
// default command executor impl
// performs http.Client request
// serves as command executor middleware for all command
func (cl Client) Execute(req *http.Request) (*http.Response, error) {
	req.Header.Add("Accept", "application/json")
	return cl.HTTPClient.Do(req)
}

func (c Client) ExecuteCommandStrategy(cmd *Command, st ...CommandExecutor) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%s%s", c.BaseURL, c.Session.Route, c.Session.Id, cmd.Path)

	rr := io.LimitReader(ReusableReader(bytes.NewReader(cmd.Data)), c.RequestReaderLimit)
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

// ExecuteCommand
//  1. general purpose client receiver
//     executes prepared command and strategies (if defined)
//     when no strategy difened, executes client request
//  2. returns response wrapper for multiple reads
func (c *Client) ExecuteCommand(cmd *Command) []*buffResponse {
	req, err := newCommandRequest(c, cmd)
	if err != nil {
		return nil
	}

	st := newExecutorContext(c, cmd)
	for i, s := range st.cmds {
		res, err := NewStrategy(s).Exec(req)
		if err != nil {
			return nil
		}

		st.bufs[i] = newBuffResponse(res)
	}

	return st.bufs
}
