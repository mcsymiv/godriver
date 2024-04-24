package driver

import (
	"net/http"

	"github.com/mcsymiv/godriver/by"
)

// newFindCommand
// returns default values for
// /element command to execute
func newFindCommand(by by.Selector, d *Driver) *Command {
	return &Command{
		Path:   "/element",
		Method: http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategies: []CommandExecutor{
			&findStrategy{
				driver:  *d,
				timeout: 20, // in 15 seconds time window performs up to 2 retries to find element
				delay:   700,
			},
		},
	}
}

func find(by by.Selector, d *Driver) (*Element, error) {
	op := newFindCommand(by, d)

	el := new(struct{ Value map[string]string })
	d.Client.ExecuteCmd(op, el)
	eId := elementID(el.Value)

	return &Element{
		Id:       eId,
		Driver:   d,
		Selector: by,
	}, nil
}

func finds(by by.Selector, d *Driver) ([]*Element, error) {
	op := &Command{
		Path:   "/elements",
		Method: http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategies: []CommandExecutor{
			&findStrategy{
				driver:  *d,
				timeout: 15, // in 15 seconds time window performs up to 2 retries to find element
				delay:   700,
			},
		},
	}

	el := new(struct{ Value []map[string]string })
	d.Client.ExecuteCmd(op, el)
	elementsId := elementsID(el.Value)

	var els []*Element

	for _, id := range elementsId {
		els = append(els, &Element{
			Id:     id,
			Driver: d,
		})
	}

	return els, nil
}

func (e *Element) From(by by.Selector) *Element {
	op := &Command{
		Path:           "/element/%s/element",
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategies: []CommandExecutor{
			&findStrategy{
				driver:  *e.Driver,
				timeout: 20, // in 15 seconds time window performs up to 2 retries to find element
				delay:   700,
			},
		},
	}

	el := new(struct{ Value map[string]string })
	e.Driver.Client.ExecuteCmd(op, el)
	eId := elementID(el.Value)

	return &Element{
		Id:       eId,
		Driver:   e.Driver,
		Selector: by,
	}
}
