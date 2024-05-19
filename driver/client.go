package driver

import (
	"bytes"
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
		RequestReaderLimit: 4096,

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

func (c *Client) ExecuteCommand(cmd *Command, d any) {
	execCtx := &execContext{
		cmd: c,
	}

	if cmd.Strategy != nil {
		execCtx.cmd = cmd.Strategy
	}

	execCtx.cmd.exec(cmd, d)
}

func (cl *Client) exec(cmd *Command, any interface{}) {
	var cPath string = cmd.Path
	if len(cmd.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(cmd.Path, cmd.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", cl.BaseURL, cPath)

	req, err := http.NewRequest(cmd.Method, url, bytes.NewBuffer(cmd.Data))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := cl.HTTPClient.Do(req)
	if err != nil {
		log.Println("error on strategy exec:", err)
		panic(err)
	}

	if any != nil {
		err = json.NewDecoder(res.Body).Decode(any)
		if err != nil {
			log.Println("error on strategy exec:", err)
			panic(err)
		}
	}

	defer res.Body.Close()
}
