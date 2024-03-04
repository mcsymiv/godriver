package driver

import (
	"log"
	"net/http"
)

func (d *Driver) GetActive() *Element {
	el, err := getActive(d)
	if err != nil {
		log.Println("error on get active", err)
		return nil
	}

	return el
}

func getActive(d *Driver) (*Element, error) {
	op := &Command{
		Path:   "/element/active",
		Method: http.MethodGet,
	}

	bRes := d.Client.ExecuteCommand(op)
	el := new(struct{ Value map[string]string })
	unmarshalResponses(bRes, el)
	eId := elementID(el.Value)

	return &Element{
		Id:     eId,
		Driver: d,
	}, nil
}
