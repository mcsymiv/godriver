package driver

import (
	"log"
	"net/http"
)

type clickStrategy struct {
	http.Client
}

func (cl clickStrategy) Execute(req *http.Request) (*http.Response, error) {
	log.Printf("click on: %s", req.URL.Path)
	return cl.Client.Do(req)
}

func (el *Element) click() (*Element, error) {
	op := &Command{
		Path:           "/element/%s/click",
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodPost,
		Data:           marshalData(&Empty{}),
		Strategies: []CommandExecutor{
			&clickStrategy{},
		},
	}

	el.Client.ExecuteCommand(op)

	return el, nil
}

func (el *Element) Click() *Element {
	e, err := el.click()
	if err != nil {
		return nil
	}

	return e
}
