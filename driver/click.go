package driver

import (
	"net/http"
)

func (el *Element) Click() *Element {
	el.Client.ExecuteCmd(&Command{
		Path:           "/element/%s/click",
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodPost,
		Data:           marshalData(&Empty{}),
	})

	return el
}

func (el *Element) Cl() *Element {
	_, err := el.Client.ExecuteCmd(&Command{
		Path:           "/element/%s/click",
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodPost,
		Data:           marshalData(&Empty{}),
	})

	if err != nil {
		return nil
	}

	return el
}
