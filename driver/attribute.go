package driver

import (
	"net/http"
)

func (e *Element) Attr(a string) string {
	attrResponse := new(struct{ Value string })

	cmd := attrCommand(e, a)
	cmd.Strategy = &attrStrategy{e.Driver}

	e.Client.ExecuteCommand(cmd, attrResponse)

	return attrResponse.Value
}

func (e *Element) HasAttr(a string) bool {
	var hasAttr bool

	e.Client.ExecuteCommand(&Command{
		PathFormatArgs: []any{e.Id, a},
		Path:           PathElementAttribute,
		Method:         http.MethodGet,
		Strategy: &hasAttributeStrategy{
			Driver:        e.Driver,
			attrToContain: a,
		},
	}, hasAttr)

	return hasAttr
}

func (e *Element) IsAttr(a string) bool {
	attrResponse := new(struct{ Value string })

	cmd := attrCommand(e, a)
	cmd.Strategy = &attrStrategy{e.Driver}

	e.Client.ExecuteCommand(cmd, attrResponse)

	return attrResponse.Value == a
}

func attrCommand(e *Element, attr string) *Command {
	return &Command{
		PathFormatArgs: []any{e.Id, attr},
		Path:           PathElementAttribute,
		Method:         http.MethodGet,
	}
}
