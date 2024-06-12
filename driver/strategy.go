package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mcsymiv/godriver/config"
)

// StrategyExecutor
type StrategyExecutor interface {
	execute(Driver)
}

type retryRequester interface {
	verify(*http.Response, interface{}) bool
}

func (d Driver) execute(st StrategyExecutor) {
	st.execute(d)
}

// strategy for strategy
// perform request command in loop
// until TimeoutFind is reached
// TODO: allow changes to TestSettings in strategy
// or until "true" breaks command loop
type retryStrategyRequest struct {
	retryRequester
}

// defaultStrategy
// executes simple Command Request
type defaultStrategy struct {
	Command
}

// retryStrategy
// executes Command Request
// N times until config.TestSetting.TimeoutFind is reached
// TimeoutFind, i.e. find elemenent timeout original naming
// Usage for commands that require explicit await set globaly
type retryStrategy struct {
	Command
}

type displayStrategy struct {
	Command
}

type isDisplayStrategy struct {
	Command
}

type hasAttributeStrategy struct {
	Command
	attrToContain string
}

func (f retryStrategy) execute(d Driver) {
	r := retryStrategyRequest{f}
	r.perform(f.Command, d)
}

func (f displayStrategy) execute(d Driver) {
	v := retryStrategyRequest{f}
	v.perform(f.Command, d)
}

func (f isDisplayStrategy) execute(d Driver) {
	v := retryStrategyRequest{f}
	v.perform(f.Command, d)
}

func (f hasAttributeStrategy) execute(d Driver) {
	v := retryStrategyRequest{f}
	v.perform(f.Command, d)
}

// perform
// wraps loopRequest strategy
// unravels cmd data, i.e. post body, url etc.
// performs NewRequest in loop
func (r retryStrategyRequest) perform(cmd Command, d Driver) {
	var cPath string = cmd.Path
	if len(cmd.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(cmd.Path, cmd.PathFormatArgs...)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := fmt.Sprintf("%s%s", d.Client.BaseURL, cPath)
	start := time.Now()
	end := start.Add(config.TestSetting.TimeoutFind * time.Second)

	for {
		req, err := http.NewRequestWithContext(ctx, cmd.Method, url, bytes.NewBuffer(cmd.Data))
		if err != nil {
			log.Println("error on NewRequest")
			panic(err)
		}

		res, err := d.Client.HTTPClient.Do(req)
		if err != nil {
			log.Println("error on Client Do Request")
			res.Body.Close()
			panic(err)
		}

		// strategy for strategy
		// "verified" response will return true
		// and break out of the loop
		if r.verify(res, cmd.ResponseData) {
			break
		}

		// close res res.Body if not verified
		// i.e. loopStrategyRequest returns false
		res.Body.Close()

		if time.Now().After(end) {
			if config.TestSetting.ScreenshotOnFail {
				d.Screenshot()
			}

			break
		}

		time.Sleep(config.TestSetting.TimeoutDelay * time.Millisecond)
		log.Println("retry find element")
	}
}

func (f retryStrategy) verify(res *http.Response, b interface{}) bool {
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

func (d displayStrategy) verify(res *http.Response, b interface{}) bool {
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
func (a isDisplayStrategy) verify(res *http.Response, b interface{}) bool {
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

func (h hasAttributeStrategy) verify(res *http.Response, a interface{}) bool {
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

func (d defaultStrategy) execute(dr Driver) {
	var cPath string = d.Command.Path
	if len(d.Command.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(d.Command.Path, d.Command.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", dr.Client.BaseURL, cPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, d.Command.Method, url, bytes.NewBuffer(d.Command.Data))
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
