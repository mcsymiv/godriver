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

type findStrategy struct {
	driver *Driver
}

func newFindStrategy(d *Driver) *findStrategy {
	return &findStrategy{
		driver: d,
	}
}

func (f *findStrategy) exec(cmd *Command, any interface{}) {
	var cPath string = cmd.Path
	if len(cmd.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(cmd.Path, cmd.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", f.driver.Client.BaseURL, cPath)
	start := time.Now()
	end := start.Add(config.TestSetting.TimeoutFind * time.Second)

	for {
		req, err := http.NewRequest(cmd.Method, url, bytes.NewBuffer(cmd.Data))
		if err != nil {
			log.Println("error on NewRequest")
			panic(err)
		}

		res, err := f.driver.Client.HTTPClient.Do(req)
		if err != nil {
			log.Println("error on Client Do Request")
			res.Body.Close()
			panic(err)
		}

		if res.StatusCode == http.StatusOK {
			err = json.NewDecoder(res.Body).Decode(any)
			if err != nil {
				log.Println("error on json NewDecoder")
				panic(err)
			}

			res.Body.Close()
			break
		}

		res.Body.Close()

		if time.Now().After(end) {

			if config.TestSetting.ScreenshotOnFail {
				f.driver.Screenshot()
			}

			break
		}

		log.Println("retry find element")
	}
}

type displayStrategy struct {
	*Driver
}

func (f *displayStrategy) exec(cmd *Command, any interface{}) {
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
