package driver

import (
	"fmt"
	"net/http"
)

// IsDisplayed
func (e *Element) IsDisplayed() *Element {
	dis, err := isDisplayed(e)
	if err != nil {
		return nil
	}

	if !dis {
		return nil
	}

	return e
}

// Is("displayed")
func (e *Element) Is() *Element {
	dis, err := isDisplayed(e)
	if err != nil {
		return nil
	}

	if !dis {
		return nil
	}

	return e
}

// newDisplayCommand
// referrences find elememnt cmd
func newDisplayCommand(e *Element) *Command {
	fCommand := newFindCommand(e.Selector, e.Driver)
	fReq, _ := newCommandRequest(e.Client, fCommand)

	return &Command{
		Path:           "/element/%s/displayed",
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodGet,

		Strategies: []CommandExecutor{
			&displayStrategy{
				Driver:  e.Driver,
				findReq: fReq,
				timeout: 8,
				delay:   800,
			},
		},
	}
}

func isDisplayed(e *Element) (bool, error) {
	op := newDisplayCommand(e)

	d := new(struct{ Value bool })
	_, err := e.Client.ExecuteCmd(op, d)
	if err != nil {
		return false, fmt.Errorf("error on display: %v", err)
	}

	return d.Value, nil
}
