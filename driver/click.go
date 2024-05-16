package driver

import (
	"fmt"
	"net/http"
)

func (el *Element) Click() *Element {

	if el.ElementError != nil {
		return &Element{
			ElementError: fmt.Errorf("element to click got error: %v", el.ElementError),
		}
	}

	bRes, err := el.Client.ExecuteCmd(&Command{
		Path:           PathElementClick,
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodPost,
		Data:           marshalData(&Empty{}),
	})

	if err != nil {
		printBuffRes(bRes)
		return &Element{
			ElementError: fmt.Errorf("unable to click on element: %v, got: %v", el, err),
		}
	}

	return el
}
