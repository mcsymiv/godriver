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
		Client:  &http.Client{},
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

	fmt.Println("inside find stratedy")
	res, err = f.Client.Do(req)

	if res.StatusCode == http.StatusNotFound {
		log.Printf("element not fount: %v", res.StatusCode)

		start := time.Now()
		end := start.Add(f.timeout * time.Second)

		for {
			log.Println("find retry")
			time.Sleep(f.delay * time.Millisecond)
			res, err = f.Client.Do(req)
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
	http.Client
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
	var res *http.Response
	var err error

	// perform isDisplayed check
	res, err = dis.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// get response buffer
	// reads response body
	buffRes, err = newBuffResponse(res)
	if err != nil {
		return nil, fmt.Errorf("error on isDisplay strategy, new buffered response: %v", err)
	}

	unmarshalResponses([]*buffResponse{buffRes}, displayRes)

	// start waiter check
	if !displayRes.Value {
		start := time.Now()
		end := start.Add(dis.timeout * time.Second)
		log.Printf("element is not visible: %v", displayRes)

		for {
			time.Sleep(dis.delay * time.Millisecond)

			res, err = dis.Client.HTTPClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("error")
			}

			buffRes, err = newBuffResponse(res)
			if err != nil {
				return nil, fmt.Errorf("error on isDisplay value retry, new buffered response: %v", err)
			}
			unmarshalResponses([]*buffResponse{buffRes}, displayRes)

			if displayRes.Value {
				// set NopCloser response with body
				buffRes.Response.Body = buffRes.bRead()
				return buffRes.Response, nil
			}

			if time.Now().After(end) {
				dis.Screenshot()
				return buffRes.Response, fmt.Errorf("error on element display timeout: %v", err)
			}
		}
	}

	// set NopCloser response with body
	buffRes.Response.Body = buffRes.bRead()
	return buffRes.Response, err
}
