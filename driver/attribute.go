package driver

import "net/http"

func (e *Element) Attribute(attr string) string {
	a := new(struct{ Value string })
	e.Client.ExecuteCmd(&Command{
		PathFormatArgs: []any{e.Id, attr},
		Path:           "/element/%s/attribute/%s",
		Method:         http.MethodGet,
	}, a)

	return a.Value
}
