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
			log.Println("in loop")
			time.Sleep(f.delay * time.Millisecond)
			res, err = f.Client.HTTPClient.Do(req)
			if err != nil {
				return res, fmt.Errorf("error on find retry: %v", err)
			}

			log.Println("make request to find element", res.StatusCode)

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

func (d *Driver) Find(by by.Selector) *Element {
	el, err := find(by, d)
	if err != nil {
		return nil
	}
	return el
}

func (d *Driver) FindX(selector string) *Element {
	by := by.Selector{
		Using: by.ByXPath,
		Value: selector,
	}

	el, err := find(by, d)
	if err != nil {
		return nil
	}
	return el
}

func (d *Driver) FindXs(selector string) []*Element {
	by := by.Selector{
		Using: by.ByXPath,
		Value: selector,
	}

	el, err := finds(by, d)
	if err != nil {
		return nil
	}
	return el
}

func (d *Driver) FindCss(selector string) *Element {
	by := by.Selector{
		Using: by.ByCssSelector,
		Value: selector,
	}

	el, err := find(by, d)
	if err != nil {
		log.Println(err)
		return nil
	}
	return el
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
				timeout: 15, // in 15 seconds time window performs up to 2 retries to find element
				delay:   700,
			},
		},
	}
}

func find(by by.Selector, d *Driver) (*Element, error) {
	op := newFindCommand(by, d)

	bRes := d.Client.ExecuteCommand(op)
	defer bRes[0].Response.Body.Close()

	el := new(struct{ Value map[string]string })
	unmarshalData(bRes[0].Response, &el)
	eId := elementID(el.Value)

	return &Element{
		Id:       eId,
		Driver:   d,
		Selector: by,
	}, nil
}

func finds(by by.Selector, d *Driver) ([]*Element, error) {
	op := d.Commands["finds"]
	op.Data = marshalData(&JsonFindUsing{
		Using: by.Using,
		Value: by.Value,
	})

	st := &findStrategy{
		Driver:  *d,
		timeout: 15,
		delay:   1000,
	}

	res, err := d.Client.ExecuteCommandStrategy(op, st)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	el := new(struct{ Value []map[string]string })
	unmarshalData(res, &el)
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
