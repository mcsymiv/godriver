package driver

import (
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
func (cl *Client) Execute(req *http.Request) (*http.Response, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := cl.HTTPClient.Do(req)
	if err != nil {
		log.Println("error on strategy exec:", err)
		return nil, err
	}

	return res, nil
}

func (cl *Client) Exec(r *buffRequest) (*buffResponse, error) {
	return nil, nil
}
