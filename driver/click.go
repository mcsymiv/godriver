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
		Strategies: []CommandExecutor{
			&clickStrategy{}, // if no initialized Driver provided, http.Client wrap can be used in strategy
		},
	})

	return el
}
