package driver

import (
	"fmt"
	"net/http"

	"github.com/mcsymiv/godriver/by"
)

func (el *Element) Click() *Element {

	el.Client.ExecuteCommand(&Command{
		Path:           PathElementClick,
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodPost,
		Data:           marshalData(&Empty{}),
	}, &Empty{})

	return el
}

func (d *Driver) Cl(selector string) *Element {
	w3cBy := by.Strategy(selector)

	el, err := f(w3cBy, d)
	if err != nil {
		panic(fmt.Errorf("unable to find element, got: %v\n", err))
	}

	return el.Click()
}
