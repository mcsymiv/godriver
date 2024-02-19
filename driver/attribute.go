package driver

import "net/http"

func (e *Element) Attribute(attr string) string {
	a, err := attribute(e, attr)
	if err != nil {
		return ""
	}

	return a
}

// Attribute
// Returns elements attribute value
func attribute(e *Element, a string) (string, error) {
	op := &Command{
		PathFormatArgs: []any{e.Id, a},
		Path:           "/element/%s/attribute/%s",
		Method:         http.MethodGet,
	}

	bRes := e.Client.ExecuteCommand(op)
	attr := new(struct{ Value string })
	unmarshalResponses(bRes, attr)

	return attr.Value, nil
}
