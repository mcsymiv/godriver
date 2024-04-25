package driver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

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
	res, err := l.next.RoundTrip(r)
	if err != nil {
		return nil, fmt.Errorf("error on %v request: %v", r, err)
	}

	log.Printf("Response: %v %v", r.URL.String(), res.StatusCode)
	return res, nil
}
