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
	st := displayStrategy{
		Driver: e.Driver,
		Command: Command{
			Path:           PathElementDisplayed,
			Method:         http.MethodGet,
			PathFormatArgs: []any{e.Id},
			Data:           marshalData(Empty{}),
		},
	}

	st.execute()

	return e
}

// IsDisplayed
// returns bool result after TimeoutFind timeout
// can seat in between element commands
// e.g.: d.F("selector").IsDisplayed()
func (e Element) IsDisplayed() bool {
	var is bool

	st := isDisplayStrategy{
		Driver: e.Driver,
		Command: Command{
			Path:           PathElementDisplayed,
			Method:         http.MethodGet,
			PathFormatArgs: []any{e.Id},
			Data:           marshalData(Empty{}),
			ResponseData:   is,
		},
	}

	st.execute()

	return is
}
