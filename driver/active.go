package driver

import (
	"net/http"
)

func (d *Driver) Active() *Element {
	el := new(struct{ Value map[string]string })

	d.execute(defaultStrategy{Command{
		Path:         PathElementActive,
		Method:       http.MethodGet,
		ResponseData: el,
	}})

	eId := elementID(el.Value)

	return &Element{
		Id:     eId,
		Driver: d,
	}
}
