package driver

import (
	"log"
	"net/http"
)

// IsDisplayed
func (e *Element) IsDisplayed() *Element {
	dis, err := isDisplayed(e)
	if err != nil {
		log.Println("error on displayed")
		return nil
	}

	if !dis {
		log.Println("element not visible")
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
	e.Client.ExecuteCmd(op, d)

	return d.Value, nil
}
