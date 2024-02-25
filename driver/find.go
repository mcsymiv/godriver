package driver

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mcsymiv/godriver/by"
)

type findStrategy struct {
	Driver
	timeout, delay time.Duration
}

// Execute
// findStrategy impl
// retries find command with delay until element is returned
// or timeout reached, which takes a screenshot of the page
func (f *findStrategy) Execute(req *http.Request) (*http.Response, error) {
	log.Printf("find strategy request: %s", req.URL.Path)
	var res *http.Response
	var err error

	res, err = f.Client.HTTPClient.Do(req)

	if res.StatusCode == http.StatusNotFound {
		log.Printf("element not fount: %v", res.StatusCode)

		start := time.Now()
		end := start.Add(f.timeout * time.Second)

		for {
			log.Println("find retry")
			time.Sleep(f.delay * time.Millisecond)
			res, err = f.Client.HTTPClient.Do(req)
			if err != nil {
				return res, fmt.Errorf("error on find retry: %v", err)
			}

			if res.StatusCode == http.StatusOK {
				log.Printf("element fount: %v", res.StatusCode)

				return res, nil
			}

			if time.Now().After(end) {
				f.Driver.Screenshot()
				return res, fmt.Errorf("unable to find element with %v timeout: %v", f.timeout, err)
			}
		}
	}

	log.Printf("find response strategy status code: %v", res.StatusCode)

	return res, err
}

// newFindCommand
// returns default values for
// /element command to execute
func newFindCommand(by by.Selector, d *Driver) *Command {
	return &Command{
		Path:   "/element",
		Method: http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategies: []CommandExecutor{
			&findStrategy{
				Driver:  *d,
				timeout: 20, // in 15 seconds time window performs up to 2 retries to find element
				delay:   700,
			},
		},
	}
}

func find(by by.Selector, d *Driver) (*Element, error) {
	op := newFindCommand(by, d)

	el := new(struct{ Value map[string]string })
	d.Client.ExecuteCmd(op, el)
	eId := elementID(el.Value)

	return &Element{
		Id:       eId,
		Driver:   d,
		Selector: by,
	}, nil
}

func finds(by by.Selector, d *Driver) ([]*Element, error) {
	op := &Command{
		Path:   "/elements",
		Method: http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategies: []CommandExecutor{
			&findStrategy{
				Driver:  *d,
				timeout: 15, // in 15 seconds time window performs up to 2 retries to find element
				delay:   700,
			},
		},
	}

	el := new(struct{ Value []map[string]string })
	d.Client.ExecuteCmd(op, el)
	elementsId := elementsID(el.Value)

	var els []*Element

	for _, id := range elementsId {
		els = append(els, &Element{
			Id:     id,
			Driver: d,
		})
	}

	return els, nil
}
