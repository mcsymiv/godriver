package driver

import (
	"net/http"
)

type clickStrategy struct {
	d Driver
}

// click strategy
// note: return default client click request as example
// TODO: add strategy for ElementNotFound, ClickIntercepted etc
func (cl clickStrategy) Execute(req *http.Request) (*http.Response, error) {
	return cl.d.Client.HTTPClient.Do(req)
}

func (el *Element) Click() *Element {
	el.Client.ExecuteCmd(&Command{
		Path:           "/element/%s/click",
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodPost,
		Data:           marshalData(&Empty{}),
		Strategies: []CommandExecutor{
			&clickStrategy{},
		},
	})

	return el
}
