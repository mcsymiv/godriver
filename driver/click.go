package driver

import (
	"fmt"
	"net/http"

	"github.com/mcsymiv/godriver/by"
)

func (el Element) Click() Element {
	el.Driver.execute(defaultStrategy{Command{
		Path:           PathElementClick,
		Method:         http.MethodPost,
		PathFormatArgs: []any{el.Id},
		Data:           marshalData(Empty{}),
	}})

	return el
}

// https://github.com/w3c/webdriver/issues/915#issuecomment-301100300
func (el Element) DoubleClick() Element {
	el.Driver.execute(defaultStrategy{Command{
			Path:           PathElementClick,
			Method:         http.MethodPost,
			PathFormatArgs: []any{el.Id},
			Data:           marshalData(Empty{}),
		}})

	return el
}

func (d Driver) Cl(selector string) Element {
	w3cBy := by.Strategy(selector)

	el, err := f(w3cBy, d)
	if err != nil {
		panic(fmt.Errorf("unable to find element, got: %v\n", err))
	}

	return el.Click()
}
