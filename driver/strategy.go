package driver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type findStrategy struct {
	driver *Driver
	*http.Client
	timeout, delay time.Duration
}

func newFindStrategy(d *Driver) *findStrategy {
	return &findStrategy{
		driver:  d,
		timeout: 20, // in 20 seconds time window performs up to 2 retries to find element
		delay:   700,
	}
}

func (cl findStrategy) Exec(r *buffRequest) (*buffResponse, error) {
	return nil, nil
}

// Execute
// findStrategy impl
// retries find command with delay until element is returned
// or timeout reached, which takes a screenshot of the page
func (f *findStrategy) Execute(req *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error

	res, err = f.driver.Client.HTTPClient.Do(req)

	if res.StatusCode == http.StatusNotFound {
		log.Printf("element not fount: %v", res.StatusCode)

		start := time.Now()
		end := start.Add(f.timeout * time.Second)

		for {
			log.Println("find retry")
			time.Sleep(f.delay * time.Millisecond)
			res, err = f.driver.Client.HTTPClient.Do(req)
			if err != nil {
				return res, fmt.Errorf("error on find retry: %v", err)
			}

			if res.StatusCode == http.StatusOK {
				log.Printf("element fount: %v", res.StatusCode)

				return res, nil
			}

			if time.Now().After(end) {
				f.driver.Screenshot()
				return res, fmt.Errorf("unable to find element with %v timeout: %v", f.timeout, err)
			}
		}
	}

	return res, err
}

type clickStrategy struct {
	*http.Client
}

// click strategy
// note: return default client click request as example
// TODO: add strategy for ElementNotFound, ClickIntercepted etc
func (cl clickStrategy) Execute(req *http.Request) (*http.Response, error) {
	return cl.Client.Do(req)
}

func (cl clickStrategy) Exec(r *buffRequest) (*buffResponse, error) {
	return nil, nil
}

type displayStrategy struct {
	findReq *http.Request
	*Driver
	timeout, delay time.Duration
}

func (dis displayStrategy) Exec(r *buffRequest) (*buffResponse, error) {
	return nil, nil
}

func (dis displayStrategy) Execute(req *http.Request) (*http.Response, error) {
	var displayRes = new(struct{ Value bool })
	var buffRes *buffResponse
	var err error

	// start waiter check
	for start := time.Now(); time.Since(start) < dis.timeout*time.Second; {
		res, err := dis.Driver.Client.HTTPClient.Do(req)
		if err != nil {
			err = fmt.Errorf("error on display strategy request")
			break
		}
		defer res.Body.Close()

		buffRes, err = newBuffResponse(res)
		if err != nil {
			err = fmt.Errorf("error on isDisplay value retry, new buffered response: %v", err)
			break
		}
		unmarshalResponses([]*buffResponse{buffRes}, displayRes)

		if displayRes.Value {
			break
		}

		time.Sleep(dis.delay * time.Millisecond)
	}

	buffRes.Response.Body = buffRes.bRead()

	// if not displayed, dis.Screenshot()
	if !displayRes.Value {
		dis.Screenshot()
		return buffRes.Response, fmt.Errorf("element is not displayed")
	}

	// set NopCloser response with body
	return buffRes.Response, err
}
