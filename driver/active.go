package driver

import (
	"net/http"
)

func (d *Driver) Active() *Element {
	el := new(struct{ Value map[string]string })
	_, err := d.Client.ExecuteCmd(&Command{
		Path:   "/element/active",
		Method: http.MethodGet,
	}, el)

	if err != nil {
		return nil
	}

	eId := elementID(el.Value)

	return &Element{
		Id:     eId,
		Driver: d,
	}
}
