package driver

import "net/http"

func (el *Element) Text() string {
	t := new(struct{ Value string })
	el.Client.ExecuteCommand(&Command{
		Path:           PathElementText, 
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodGet,
	}, t)

	return t.Value
}
