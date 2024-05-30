package driver

import (
	"net/http"
)

func (d *Driver) Active() *Element {
	el := new(struct{ Value map[string]string })
	d.Client.ExecuteCommand(&Command{
		Path:   PathElementActive,
		Method: http.MethodGet,
	}, el)

	eId := elementID(el.Value)

	return &Element{
		Id:     eId,
		Driver: d,
	}
}

func (d Driver) IsActive() Element {
	el := new(struct{ Value map[string]string })

	st := Strategy{
		Command: Command{
			Path:         PathElementActive,
			Method:       http.MethodGet,
			ResponseData: el,
		},
	}

	d.execute(st)

	eId := elementID(el.Value)

	return Element{
		Id:     eId,
		Driver: d,
	}
}
