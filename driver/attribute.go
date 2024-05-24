package driver

import "net/http"

func (e *Element) Attr(attr string) string {
	a := new(struct{ Value string })

	e.Client.ExecuteCommand(&Command{
		PathFormatArgs: []any{e.Id, attr},
		Path:           PathElementAttribute,
		Method:         http.MethodGet,
	}, a)

	return a.Value
}
