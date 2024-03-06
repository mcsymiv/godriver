package driver

import (
	"net/http"
)

func (d *Driver) GetActive() *Element {
	el := new(struct{ Value map[string]string })
	d.Client.ExecuteCmd(&Command{
		Path:   "/element/active",
		Method: http.MethodGet,
	}, el)
	eId := elementID(el.Value)

	return &Element{
		Id:     eId,
		Driver: d,
	}
}
