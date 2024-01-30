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
	log.Printf("before click strategy driver request: %s", req.URL.Path)
	res, err := cl.Client.Do(req)
	log.Printf("after click strategy driver request: %s", req.URL.Path)

	return res, err
}

func (el Element) Click() error {
	op := &Command{
		Path:   fmt.Sprintf("/element/%s/click", el.Id),
		Method: http.MethodPost,
		Data:   marshalData(&Empty{}),
	}

	st := &clickStrategy{}

	_, err := el.Client.ExecuteCommandStrategy(op, st)
	if err != nil {
		log.Println("error on click:", err)
		return err
	}

	return nil
}
