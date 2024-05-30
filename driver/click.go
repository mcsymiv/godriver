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

// https://github.com/w3c/webdriver/issues/915#issuecomment-301100300
func (el Element) DoubleClick() Element {
	st := Strategy{
		Command: Command{
			Path:           PathElementClick,
			Method:         http.MethodPost,
			PathFormatArgs: []any{el.Id},
			Data:           marshalData(Empty{}),
		},
	}

	el.Driver.execute(st)

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
