package driver

import (
	"net/http"
)

func (d Driver) Active() Element {
	el := new(struct{ Value map[string]string })

	st := defaultStrategy{
		Driver: d,
		Command: Command{
			Path:         PathElementActive,
			Method:       http.MethodGet,
			ResponseData: el,
		},
	}

	st.execute()

	eId := elementID(el.Value)

	return Element{
		Id:     eId,
		Driver: d,
	}
}
