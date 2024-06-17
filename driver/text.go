package driver

import "net/http"

func (el *Element) Text() string {
	t := new(struct{ Value string })

	el.Driver.execute(defaultStrategy{Command{
		Path:           PathElementText,
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodGet,
		ResponseData:   t,
	}})

	return t.Value
}
