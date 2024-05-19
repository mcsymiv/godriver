package driver

import (
	"fmt"
	"net/http"
)

// Is("displayed")
// returns chained Element
// can seat in between element commands
// e.g.: d.F("selector").Is().Attr("href")
// will panic if found elemet is not displayed
func (e *Element) Is() *Element {
	err := isDisplayed(e)
	if err != nil {
		panic(fmt.Errorf("got: %v\n", err))
	}

	return e
}

// newDisplayCommand
// referrences find elememnt cmd
func newDisplayCommand(e *Element) *Command {
	return &Command{
		Path:           PathElementDisplayed,
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodGet,

		Strategies: []CommandExecutor{
			&displayStrategy{
				Driver: e.Driver,
			},
		},
	}
}

func isDisplayed(e *Element) error {
	op := newDisplayCommand(e)

	_, err := e.Client.ExecuteCmd(op)
	if err != nil {
		return fmt.Errorf("error on display: %v\n", err)
	}

	return nil
}
