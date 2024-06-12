package driver

import (
	"net/http"
	"time"
)

// RestClient represents a REST client configuration.
type Client struct {
	BaseURL            string
	HTTPClient         *http.Client
	RequestReaderLimit int64
	// syncMutex  sync.Mutex // Mutex for ensuring thread safety
}

// newClientV2
// new client init without Session param
func newClient(baseURL string) Client {
	return Client{
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
