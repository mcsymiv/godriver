package driver

import (
	"fmt"
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
		Path:   fmt.Sprintf("/element/%s/click", el.Id),
		Method: http.MethodPost,
		Data:   marshalData(&Empty{}),
	}

	st := &clickStrategy{}

	_, err := el.Client.ExecuteCommandStrategy(op, st)
	if err != nil {
		log.Println("error on click:", err)
		return nil, err
	}

	return el, nil
}

func (el *Element) Click() *Element {
	e, err := el.click()
	if err != nil {
		return nil
	}

	return e
}
