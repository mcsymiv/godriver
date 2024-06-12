package driver

import (
	"net/http"
)

// Is("displayed")
// returns chained Element
// can seat in between element commands
// e.g.: d.F("selector").Is().Attr("href")
// will panic if found elemet is not displayed
func (e Element) Is() Element {
	e.Driver.execute(displayStrategy{Command{
		Path:           PathElementDisplayed,
		Method:         http.MethodGet,
		PathFormatArgs: []any{e.Id},
		Data:           marshalData(Empty{}),
	}})

	return e
}

// IsDisplayed
// returns bool result after TimeoutFind timeout
// can seat in between element commands
// e.g.: d.F("selector").IsDisplayed()
func (e Element) IsDisplayed() bool {
	var is bool

	e.Driver.execute(isDisplayStrategy{Command{
		Path:           PathElementDisplayed,
		Method:         http.MethodGet,
		PathFormatArgs: []any{e.Id},
		Data:           marshalData(Empty{}),
		ResponseData:   is,
	}})

	return is
}
