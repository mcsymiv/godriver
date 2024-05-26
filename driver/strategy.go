package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mcsymiv/godriver/config"
)

// strategy for strategy
// perform request command in loop
// until TimeoutFind is reached
// TODO: allow changes to TestSettings in strategy
// or until "true" breaks command loop
type loopStrategyRequest struct {
	loopRequester
	*Command
	*Driver
	a interface{} // a, any data for response decoding
}

// newStrategy
// initializes new loopRequest
// TODO: maybe reduce number of params
func newLoopStrategy(r loopRequester, c *Command, d *Driver, a interface{}) *loopStrategyRequest {
	return &loopStrategyRequest{r, c, d, a}
}

// requester
// actual loopRequest interface
// i.e. strategy for strategy
// to verify response in loop
type loopRequester interface {
	verify(*http.Response, interface{}) bool
}

// performStrategy
// wraps loopRequest strategy
// unravels cmd data, i.e. post body, url etc.
// performs NewRequest in loop
// passes response to s*loopStrategyRequest (find, display, attribute)
func (s *loopStrategyRequest) performStrategy() {
	var cPath string = s.Path
	if len(s.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(s.Path, s.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", s.Client.BaseURL, cPath)
	start := time.Now()
	end := start.Add(config.TestSetting.TimeoutFind * time.Second)

	for {
		req, err := http.NewRequest(s.Method, url, bytes.NewBuffer(s.Data))
		if err != nil {
			log.Println("error on NewRequest")
			panic(err)
		}

		res, err := s.Client.HTTPClient.Do(req)
		if err != nil {
			log.Println("error on Client Do Request")
			res.Body.Close()
			panic(err)
		}

		// strategy for strategy
		// "verified" response will return true
		// and break out of the loop
		if s.verify(res, s.a) {
			break
		}

		res.Body.Close()
		if time.Now().After(end) {
			if config.TestSetting.ScreenshotOnFail {
				s.Screenshot()
			}

			break
		}
	}
}

// command strategies section
//
// findStrategy
type findStrategy struct {
	driver *Driver
}

func (f *findStrategy) exec(cmd *Command, any interface{}) {
	v := newLoopStrategy(f, cmd, f.driver, any)
	v.performStrategy()
}

func (f *findStrategy) verify(res *http.Response, b interface{}) bool {
	if res.StatusCode == http.StatusOK {
		err := json.NewDecoder(res.Body).Decode(b)
		if err != nil {
			log.Println("error on json NewDecoder")
			panic(err)
		}

		res.Body.Close()
		return true
	}

	return false
}

// displayStrategy
type displayStrategy struct {
	*Driver
}

func (f *displayStrategy) exec(cmd *Command, any interface{}) {
	v := newLoopStrategy(f, cmd, f.Driver, any)
	v.performStrategy()
}

func (d *displayStrategy) verify(res *http.Response, b interface{}) bool {
	if res.StatusCode == http.StatusOK {
		err := json.NewDecoder(res.Body).Decode(b)
		if err != nil {
			log.Println("error on json NewDecoder")
			res.Body.Close()
			panic(err)
		}

		if b.(struct{ Value bool }).Value {
			res.Body.Close()
			return true
		}
	}

	return false
}

// clickStrategy
// serves as an example
// that <command>Strategy can divert
// from default client request
// and request to webdriver service
// can be wrapped in custom logic
type clickStrategy struct {
	*Driver
}

func (f *clickStrategy) exec(cmd *Command, any interface{}) {
	var cPath string = cmd.Path
	if len(cmd.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(cmd.Path, cmd.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", f.Driver.Client.BaseURL, cPath)

	req, err := http.NewRequest(cmd.Method, url, bytes.NewBuffer(cmd.Data))
	if err != nil {
		log.Println("error on NewRequest")
		panic(err)
	}

	res, err := f.Driver.Client.HTTPClient.Do(req)
	if err != nil {
		log.Println("error on Client Do Request")
		panic(err)
	}

	defer res.Body.Close()
}

// attrStrategy
type attrStrategy struct {
	*Driver
}

func (at *attrStrategy) exec(cmd *Command, any interface{}) {
	v := newLoopStrategy(at, cmd, at.Driver, any)
	v.performStrategy()
}

func (a *attrStrategy) verify(res *http.Response, b interface{}) bool {
	if res.StatusCode == http.StatusOK {
		err := json.NewDecoder(res.Body).Decode(b)
		if err != nil {
			log.Println("error on json NewDecoder")
			res.Body.Close()
			panic(err)
		}

		res.Body.Close()
		return true
	}

	return false
}

// isDisplayStrategy
type isDisplayStrategy struct {
	*Driver
}

func (is *isDisplayStrategy) exec(cmd *Command, any interface{}) {
	v := newLoopStrategy(is, cmd, is.Driver, any)
	v.performStrategy()
}

func (a *isDisplayStrategy) verify(res *http.Response, b interface{}) bool {
	if res.StatusCode == http.StatusOK {
		var displayResponse = new(struct{ Value bool })

		err := json.NewDecoder(res.Body).Decode(displayResponse)
		if err != nil {
			log.Println("error on json NewDecoder")
			res.Body.Close()
			panic(err)
		}

		if displayResponse.Value {
			b = true
			res.Body.Close()
			return true
		}
	}

	return false
}
