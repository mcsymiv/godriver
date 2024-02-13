package driver

import (
	"log"
	"net/http"
	"time"
)

type findStrategy struct {
	http.Client
	Timeout, Delay time.Duration
}

func (f *findStrategy) Execute(req *http.Request) (*http.Response, error) {
	log.Printf("find strategy request: %s", req.URL.Path)
	var res *http.Response
	var err error

	res, err = f.Client.Do(req)

	if res.StatusCode == http.StatusNotFound {
		log.Printf("element not fount: %v", res.StatusCode)

		start := time.Now()
		end := start.Add(f.Timeout * time.Second)

		for {
			time.Sleep(f.Delay * time.Millisecond)
			res, err = f.Client.Do(req)

			if err != nil {
				log.Printf("element not fount: %v", res.StatusCode)
			}

			if res.StatusCode == http.StatusOK {
				log.Printf("element fount: %v", res.StatusCode)

				return res, nil
			}

			if time.Now().After(end) {
				log.Printf("timeout")
				return res, err
			}
		}
	}

	log.Printf("find response strategy status code: %v", res.StatusCode)

	return res, err
}

type By struct {
	Using, Value string
}

func (d *Driver) Find(by By) *Element {
	el, err := find(by, d)
	if err != nil {
		return nil
	}
	return el
}

func (d *Driver) FindX(selector string) *Element {
	by := By{
		Using: ByXPath,
		Value: selector,
	}

	el, err := find(by, d)
	if err != nil {
		return nil
	}
	return el
}

func (d Driver) FindXs(selector string) []*Element {
	by := By{
		Using: ByXPath,
		Value: selector,
	}

	el, err := finds(by, d)
	if err != nil {
		return nil
	}
	return el
}

func (d *Driver) FindCss(selector string) *Element {
	by := By{
		Using: ByCssSelector,
		Value: selector,
	}

	el, err := find(by, d)
	if err != nil {
		return nil
	}
	return el
}

// newFindCommand
// returns default values for
// /element command to execute
func newFindCommand(by By, d *Driver) *Command {
	return &Command{
		Path:   "/element",
		Method: http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategies: []CommandExecutor{
			&findStrategy{
				Client:  *d.Client.HTTPClient,
				Timeout: 8,
				Delay:   800,
			},
		},
	}
}

func find(by By, d *Driver) (*Element, error) {
	op := newFindCommand(by, d)

	bRes := d.Client.ExecuteCommand(op)
	defer bRes[0].Response.Body.Close()

	el := new(struct{ Value map[string]string })
	unmarshalData(bRes[0].Response, &el)
	eId := elementID(el.Value)

	return &Element{
		Id:     eId,
		Driver: d,
		By:     by,
	}, nil
}

func finds(by By, d Driver) ([]*Element, error) {
	op := d.Commands["finds"]
	op.Data = marshalData(&JsonFindUsing{
		Using: by.Using,
		Value: by.Value,
	})

	st := &findStrategy{
		Client:  *d.Client.HTTPClient,
		Timeout: 15,
		Delay:   1000,
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
			Driver: &d,
		})
	}

	return els, nil
}
