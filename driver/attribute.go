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

	st := attrStrategy{
		Driver: e.Driver,
		Command: Command{
			PathFormatArgs: []any{e.Id, a},
			Path:           PathElementAttribute,
			Method:         http.MethodGet,
			ResponseData:   attrResponse,
		},
	}

	st.execute()
	return attrResponse.Value
}

func (e Element) HasAttr(a string) bool {
	var hasAttr bool

	cmd := attrCommand(e, a)
	cmd.ResponseData = hasAttr

	st := hasAttributeStrategy{
		Driver:  e.Driver,
		Command: cmd,
	}

	st.execute()

	return hasAttr
}

func (e Element) IsAttr(a string) bool {
	attrResponse := new(struct{ Value string })

	cmd := attrCommand(e, a)
	cmd.ResponseData = attrResponse

	st := defaultStrategy{
		Driver:  e.Driver,
		Command: cmd,
	}

	st.execute()

	return attrResponse.Value == a
}
