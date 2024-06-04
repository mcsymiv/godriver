package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mcsymiv/godriver/config"
)

type StrategyExecutor interface {
	execute(Driver)
}

type loopRequesterV2 interface {
	verifyV2(*http.Response, interface{}) bool
}

// strategy for strategy
// perform request command in loop
// until TimeoutFind is reached
// TODO: allow changes to TestSettings in strategy
// or until "true" breaks command loop
type loopStrategyRequestV2 struct {
	loopRequesterV2
	Command
	Driver
}

// defaultStrategy
// executes simple Command Request
// Creates new http.Client
type defaultStrategy struct {
	Command
}

type findStrategyV2 struct {
	Command
}

type displayStrategy struct {
	Command
}

type isDisplayStrategy struct {
	Command
}

// attrStrategy
type attrStrategy struct {
	Command
}

type clickStrategy struct {
	Command
}

type hasAttributeStrategy struct {
	Command
	attrToContain string
}

func (d Driver) execute(st StrategyExecutor) {
	st.execute(d)
}

func (f findStrategyV2) execute(d Driver) {
	v := loopStrategyRequestV2{f, f.Command, d}
	v.performStrategyV2()
}

func (f displayStrategy) execute(d Driver) {
	v := loopStrategyRequestV2{f, f.Command, d}
	v.performStrategyV2()
}

func (is isDisplayStrategy) execute(d Driver) {
	v := loopStrategyRequestV2{is, is.Command, d}
	v.performStrategyV2()
}

func (is attrStrategy) execute(d Driver) {
	v := loopStrategyRequestV2{is, is.Command, d}
	v.performStrategyV2()
}

func (is hasAttributeStrategy) execute(d Driver) {
	v := loopStrategyRequestV2{is, is.Command, d}
	v.performStrategyV2()
}

func (d defaultStrategy) execute(dr Driver) {
	var cPath string = d.Command.Path
	if len(d.Command.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(d.Command.Path, d.Command.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", dr.Client.BaseURL, cPath)

	req, err := http.NewRequest(d.Command.Method, url, bytes.NewBuffer(d.Command.Data))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := dr.Client.HTTPClient.Do(req)
	if err != nil {
		log.Println("error on strategy exec:", err)
		res.Body.Close()
		panic(err)
	}

	if d.Command.ResponseData != nil {
		err = json.NewDecoder(res.Body).Decode(d.Command.ResponseData)
		if err != nil {
			log.Println("error on strategy exec:", err)
			res.Body.Close()
			panic(err)
		}
	}

	defer res.Body.Close()
}

// performStrategy
// wraps loopRequest strategy
// unravels cmd data, i.e. post body, url etc.
// performs NewRequest in loop
// passes response to s*loopStrategyRequest (find, display, attribute)
func (s loopStrategyRequestV2) performStrategyV2() {
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
		if s.verifyV2(res, s.Command.ResponseData) {
			break
		}

		// close res res.Body if not verified
		// i.e. loopStrategyRequest returns false
		res.Body.Close()

		if time.Now().After(end) {
			if config.TestSetting.ScreenshotOnFail {
				s.Screenshot()
			}

			break
		}

		time.Sleep(config.TestSetting.TimeoutDelay * time.Millisecond)
		log.Println("retry find element")
	}
}

func (f findStrategyV2) verifyV2(res *http.Response, b interface{}) bool {
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

func (d displayStrategy) verifyV2(res *http.Response, b interface{}) bool {
	if res.StatusCode == http.StatusOK {
		var displayResponse = new(struct{ Value bool })
		err := json.NewDecoder(res.Body).Decode(displayResponse)
		if err != nil {
			log.Println("error on json NewDecoder")
			res.Body.Close()
			panic(err)
		}

		if displayResponse.Value {
			res.Body.Close()
			return true
		}
	}

	return false
}

// verify isDisplayStreategy
// will assign true to b to reuse in IsDisplayed()
func (a isDisplayStrategy) verifyV2(res *http.Response, b interface{}) bool {
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

		res.Body.Close()
		return false
	}

	return false
}

func (a attrStrategy) verifyV2(res *http.Response, b interface{}) bool {
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

func (h hasAttributeStrategy) verifyV2(res *http.Response, a interface{}) bool {
	if res.StatusCode == http.StatusOK {
		var attributeResponse = new(struct{ Value string })
		err := json.NewDecoder(res.Body).Decode(attributeResponse)
		if err != nil {
			log.Println("error on json NewDecoder")
			res.Body.Close()
			panic(err)
		}

		if strings.Contains(attributeResponse.Value, h.attrToContain) {
			a = true
			res.Body.Close()
			return true
		}

		res.Body.Close()
		return false
	}

	return false
}
