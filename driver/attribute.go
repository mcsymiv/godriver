package driver

import (
	"net/http"
	"strings"
)

func (e *Element) Attr(a string) string {
	return attr(e, a)
}

func (e *Element) HasAttr(a string) bool {
	atttribute := attr(e, a)
	return strings.Contains(atttribute, a)
}

func (e *Element) IsAttr(a string) bool {
	atttribute := attr(e, a)
	return atttribute == a
}

func attr(e *Element, attr string) string {
	a := new(struct{ Value string })

	e.Client.ExecuteCommand(&Command{
		PathFormatArgs: []any{e.Id, attr},
		Path:           PathElementAttribute,
		Method:         http.MethodGet,
		Strategy: &attrStrategy{
			e.Driver,
		},
	}, a)

	return a.Value

}
