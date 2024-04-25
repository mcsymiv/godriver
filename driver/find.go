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

		Strategies: []CommandExecutor{newFindStrategy(d)},
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
		Strategies: []CommandExecutor{newFindStrategy(d)},
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
