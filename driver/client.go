package driver

import (
	"encoding/json"
	"fmt"
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
// serves as command executor middleware for all commands
func (cl Client) Execute(req *http.Request) (*http.Response, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := cl.HTTPClient.Do(req)
	if err != nil {
		log.Println("error on strategy exec:", err)
		return nil, err
	}

	return res, nil
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
			log.Println("error on strategy exec:", err)
			return nil
		}

		st.bufs[i] = newBuffResponse(res)
	}

	return st.bufs
}

// ExecuteCmd
//  1. general purpose client receiver
//     executes prepared command and strategies (if defined)
//     when no strategy difened, executes client request
//  2. unmarshals passed data struct
func (c *Client) ExecuteCmd(cmd *Command, d ...any) {
	req, err := newCommandRequest(c, cmd)
	if err != nil {
		log.Println("error on new command request")
	}

	st := newExecutorContext(c, cmd)
	for i, s := range st.cmds {
		res, err := NewStrategy(s).Exec(req)
		if err != nil {
			log.Printf("error on new strategy exec: %+v", err)
			return
		}

		st.bufs[i] = newBuffResponse(res)
	}

	if len(st.bufs) > 0 && len(d) > 0 {
		for i, res := range st.bufs {
			err := json.Unmarshal(res.buff, d[i])

			if err != nil {
				log.Printf("error on unmarshal %d response: %v", i, err)
			}
		}
	}
}
