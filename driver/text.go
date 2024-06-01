package driver

import "net/http"

func (el Element) Text() string {
	t := new(struct{ Value string })

	st := defaultStrategy{
		Driver: el.Driver,
		Command: Command{
			Path:           PathElementText,
			PathFormatArgs: []any{el.Id},
			Method:         http.MethodGet,
			ResponseData:   t,
		},
	}

	st.execute()

	return t.Value
}
