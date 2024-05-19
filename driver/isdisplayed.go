package driver

import (
	"net/http"
)

// Is("displayed")
// returns chained Element
// can seat in between element commands
// e.g.: d.F("selector").Is().Attr("href")
// will panic if found elemet is not displayed
func (e *Element) Is() *Element {
	e.Client.ExecuteCommand(&Command{
		Path:           PathElementDisplayed,
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodGet,

		Strategy: &displayStrategy{
			e.Driver,
		},
	}, nil)

	return e
}
