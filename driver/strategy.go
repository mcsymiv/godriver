package driver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mcsymiv/godriver/config"
)

type findStrategy struct {
	driver *Driver
	*http.Client

	timeout time.Duration
	delay   time.Duration
}

func newFindStrategy(d *Driver) *findStrategy {
	return &findStrategy{
		driver: d,
	}
}

func (cl findStrategy) Exec(r *buffRequest) (*buffResponse, error) {
	return nil, nil
}

func drainBody(resp *http.Response) {
	if resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func copyReqBody(req *http.Request) *http.Request {
	rUse := ReusableReader(req.Body)
	rr := io.LimitReader(rUse, 4096)
	reqBody := io.NopCloser(rr)
	// req.Body = reqBody

	req2, err := http.NewRequest(req.Method, req.URL.String(), reqBody)
	if err != nil {
		log.Println("unable to copy request: %v", err)
	}

	return req2
}

// Execute
// findStrategy impl
// retries find command with delay until element is returned
// or timeout reached, which takes a screenshot of the page
func (f *findStrategy) Execute(req *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error

	req = copyReqBody(req)
	res, err = f.driver.Client.HTTPClient.Do(req)

	if res.StatusCode == http.StatusNotFound {
		log.Printf("element not fount: %v", res.StatusCode)

		start := time.Now()
		end := start.Add(config.TestSetting.TimeoutFind * time.Second)

		for {
			log.Println("find retry")
			time.Sleep(config.TestSetting.TimeoutDelay * time.Millisecond)

			drainBody(res)
			res, err = f.driver.Client.HTTPClient.Do(req)
			if err != nil {
				log.Println("error on client request", err)
				err = fmt.Errorf("error on find retry: %v", err)
				break
			}

			if res.StatusCode == http.StatusOK {
				break
			}

			if time.Now().After(end) {

				if config.TestSetting.ScreenshotOnFail {
					f.driver.Screenshot()
				}

				err = fmt.Errorf("unable to find element with %v timeout: %v", config.TestSetting.TimeoutFind, err)
				break
			}

			if config.TestSetting.RefreshOnFindError {
				f.driver.Refresh()
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
	*Driver
}

func (dis displayStrategy) Exec(r *buffRequest) (*buffResponse, error) {
	return nil, nil
}

func (dis displayStrategy) Execute(req *http.Request) (*http.Response, error) {
	var displayRes = new(struct{ Value bool })
	var res *http.Response

	res, err := dis.Driver.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on is display request")
	}

	start := time.Now()
	end := start.Add(config.TestSetting.TimeoutFind * time.Second)

	for {
		time.Sleep(config.TestSetting.TimeoutDelay * time.Millisecond)

		if res.StatusCode == http.StatusOK {
			err = unmarshalRes(res, displayRes)
			if err != nil {
				return nil, fmt.Errorf("error on unmarshal display body, got: %v", err)
			}

			if displayRes.Value {
				return res, nil
			}
		}

		if time.Now().After(end) {
			if config.TestSetting.ScreenshotOnFail {
				dis.Screenshot()
			}

			return nil, fmt.Errorf("element not displayed within %v timeout, got: %v", config.TestSetting.TimeoutFind, err)
		}

		res, err = dis.Driver.Client.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error on display request")
		}
	}
}
