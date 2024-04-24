package driver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type findStrategy struct {
	driver         Driver
	timeout, delay time.Duration
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
