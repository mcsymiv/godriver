package driver

import (
	"net/http"
)

func attrCommand(e Element, a string) Command {
	return Command{
		PathFormatArgs: []any{e.Id, a},
		Path:           PathElementAttribute,
		Method:         http.MethodGet,
	}
}

func (e Element) Attr(a string) string {
	attrResponse := new(struct{ Value string })

	e.Driver.execute(attrStrategy{
		Command: Command{
			PathFormatArgs: []any{e.Id, a},
			Path:           PathElementAttribute,
			Method:         http.MethodGet,
			ResponseData:   attrResponse,
		},
	})

	return attrResponse.Value
}

func (e Element) HasAttr(a string) bool {
	var hasAttr bool

	cmd := attrCommand(e, a)
	cmd.ResponseData = hasAttr

	e.Driver.execute(hasAttributeStrategy{cmd, a})

	return hasAttr
}

func (e Element) IsAttr(a string) bool {
	attrResponse := new(struct{ Value string })

	cmd := attrCommand(e, a)
	cmd.ResponseData = attrResponse

	e.Driver.execute(defaultStrategy{cmd})

	return attrResponse.Value == a
}
