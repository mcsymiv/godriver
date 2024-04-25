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
	BaseURL            string
	HTTPClient         *http.Client
	RequestReaderLimit int64
	// syncMutex  sync.Mutex // Mutex for ensuring thread safety
	DefaultExecuteStrategy CommandExecutor
}

// newClientV2
// new client init without Session param
func newClient(baseURL string) *Client {
	return &Client{
		BaseURL:            baseURL,
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
// serves as command executor middleware for all default commands,
// i.e. withoud defined CommandExecutor stategy
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

// ExecuteCmd
//  1. general purpose client receiver
//     executes prepared command and strategies (if defined)
//     when no strategy difened, executes client request
//  2. unmarshals passed data struct
//
// TODO: refactor to internal use, i.e. executeCmd
func (c *Client) ExecuteCmd(cmd *Command, d ...any) ([]*buffResponse, error) {
	req, err := newCommandRequest(c, cmd)
	if err != nil {
		return nil, fmt.Errorf("error on new command request: %v", err)
	}

	st := newExecutorContext(c, cmd)
	for i, s := range st.cmds {

		// executes request inside defined CommandExecutor strategy
		// if none provided, performs http.Request with Client's DefaultExecuteStrategy
		res, err := NewStrategy(s).Exec(req)
		if err != nil {
			return nil, fmt.Errorf("error on new strategy exec: %+v", err)
		}

		st.bufs[i], err = newBuffResponse(res)
		if err != nil {
			return nil, fmt.Errorf("error on new buffered response: %v", err)
		}
	}

	if len(st.bufs) > 0 && len(d) > 0 {
		for i, res := range st.bufs {
			err := json.Unmarshal(res.buff, d[i])

			if err != nil {
				return nil, fmt.Errorf("error on unmarshal %d response: %v", i, err)
			}
		}
	}

	return st.bufs, nil
}
