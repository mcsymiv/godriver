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
// rename to loopRequest
type strategy struct {
	req requester
	cmd *Command
	dr  *Driver
	a   interface{}
}

// newStrategy
// initializes new loopRequest
// TODO: maybe reduce number of params
func newStrategy(r requester, c *Command, d *Driver, a interface{}) *strategy {
	return &strategy{
		req: r,
		cmd: c,
		dr:  d,
		a:   a,
	}
}

// requester
// actual loopRequest interface
// i.e. strategy for strategy
// to verify response in loop
type requester interface {
	// makeRequest(*http.Response, interface{}) bool // maybe change name to verifyResponse
	verify(*http.Response, interface{}) bool
}

// performStrategy
// wraps loopRequest strategy
// unravels cmd data, i.e. post body, url etc.
// performs NewRequest in loop
// passed response to s*strategy (find, display, attribute)
func (s *strategy) performStrategy() {
	var cPath string = s.cmd.Path
	if len(s.cmd.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(s.cmd.Path, s.cmd.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", s.dr.Client.BaseURL, cPath)
	start := time.Now()
	end := start.Add(config.TestSetting.TimeoutFind * time.Second)

	for {
		req, err := http.NewRequest(s.cmd.Method, url, bytes.NewBuffer(s.cmd.Data))
		if err != nil {
			log.Println("error on NewRequest")
			panic(err)
		}

		res, err := s.dr.Client.HTTPClient.Do(req)
		if err != nil {
			log.Println("error on Client Do Request")
			res.Body.Close()
			panic(err)
		}

		// strategy for strategy
		if s.req.verify(res, s.a) {
			break
		}

		res.Body.Close()
		if time.Now().After(end) {
			if config.TestSetting.ScreenshotOnFail {
				s.dr.Screenshot()
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
	v := newStrategy(f, cmd, f.driver, any)
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
	v := newStrategy(f, cmd, f.Driver, any)
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

func (a *attrStrategy) verify(res *http.Response, b interface{}) bool {
	if res.StatusCode == http.StatusOK {
		err := json.NewDecoder(res.Body).Decode(b)
		if err != nil {
			log.Println("error on json NewDecoder")
			res.Body.Close()
			panic(err)
		}

		res.Body.Close()
		return false
	}

	return false
}

func (at *attrStrategy) exec(cmd *Command, any interface{}) {
	v := newStrategy(at, cmd, at.Driver, any)
	v.performStrategy()
}

type isDisplayStrategy struct {
	*Driver
}

func (f *isDisplayStrategy) exec(cmd *Command, any interface{}) {
	var displayResponse = new(struct{ Value bool })
	var cPath string = cmd.Path
	if len(cmd.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(cmd.Path, cmd.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", f.Client.BaseURL, cPath)
	start := time.Now()
	end := start.Add(config.TestSetting.TimeoutFind * time.Second)

	for {
		req, err := http.NewRequest(cmd.Method, url, bytes.NewBuffer(cmd.Data))
		if err != nil {
			log.Println("error on NewRequest")
			panic(err)
		}

		res, err := f.Client.HTTPClient.Do(req)
		if err != nil {
			log.Println("error on Client Do Request")
			res.Body.Close()
			panic(err)
		}

		if res.StatusCode == http.StatusOK {
			err = json.NewDecoder(res.Body).Decode(displayResponse)
			if err != nil {
				log.Println("error on json NewDecoder")
				res.Body.Close()
				panic(err)
			}

			if displayResponse.Value {
				any = true
				res.Body.Close()
				break
			}
		}

		res.Body.Close()

		if time.Now().After(end) {
			if config.TestSetting.ScreenshotOnFail {
				f.Screenshot()
			}

			break
		}

		log.Println("retry find element")
	}
}
