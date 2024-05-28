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
	cmd := &Command{}
	displayCommand(e, cmd)

	cmd.Strategy = &displayStrategy{e.Driver}
	e.Client.ExecuteCommand(cmd, nil)

	return e
}

// IsDisplayed
// returns bool result after TimeoutFind timeout
// can seat in between element commands
// e.g.: d.F("selector").IsDisplayed()
func (e *Element) IsDisplayed() bool {
	var is bool

	cmd := &Command{}
	displayCommand(e, cmd)
	cmd.Strategy = &isDisplayStrategy{e.Driver}

	e.Client.ExecuteCommand(cmd, &is)

	return is
}

func displayCommand(e *Element, cmd *Command) {
	cmd.Path = PathElementDisplayed
	cmd.PathFormatArgs = []any{e.Id}
	cmd.Method = http.MethodGet
}
