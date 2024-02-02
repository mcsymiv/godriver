package driver

import (
	"log"
	"net/http"
	"time"
)

type findStrategy struct {
	http.Client
}

func (f *findStrategy) Execute(req *http.Request) (*http.Response, error) {
	log.Printf("find strategy request: %s", req.URL.Path)
	var res *http.Response
	var err error

	res, err = f.Client.Do(req)
	if res.StatusCode == http.StatusNotFound {
		log.Printf("element not fount: %v", res.StatusCode)

		start := time.Now()
		end := start.Add(10 * time.Second)
		for {
			time.Sleep(500 * time.Millisecond)
			res, err = f.Client.Do(req)
			if err != nil {
				log.Printf("element not fount: %v", res.StatusCode)
				return nil, err
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

func (d Driver) FindElement(selector string) (*Element, error) {
	jFind := &JsonFindUsing{
		Using: ByXPath,
		Value: selector,
	}

	op := &Command{
		Path:   "/element",
		Method: http.MethodPost,
		Data:   marshalData(jFind),
	}

	st := &findStrategy{
		Client: *d.Client.HTTPClient,
	}

	res, err := d.Client.ExecuteCommandStrategy(op, st)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	el := new(struct{ Value map[string]string })
	unmarshalData(res, &el)
	eId := elementID(el.Value)

	return &Element{
		Id:      eId,
		Client:  d.Client,
		Session: d.Session,
	}, nil

}
