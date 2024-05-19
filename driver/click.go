package driver

import (
	"net/http"
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
