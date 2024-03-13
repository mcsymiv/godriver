package driver

import (
	"log"
	"net/http"
)

type clickStrategy struct {
	dr *Driver
}

func (cl clickStrategy) Execute(req *http.Request) (*http.Response, error) {
	log.Printf("click on: %s", req.URL.Path)
	return cl.dr.Client.HTTPClient.Do(req)
}

func (el *Element) Click() *Element {
	el.Client.ExecuteCmd(&Command{
		Path:           "/element/%s/click",
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodPost,
		Data:           marshalData(&Empty{}),
		Strategies: []CommandExecutor{
			&clickStrategy{
				dr: el.Driver,
			},
		},
	})

	return el
}
