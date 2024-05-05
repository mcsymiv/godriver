package driver

import (
	"fmt"
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

func find(b by.Selector, d *Driver) (*Element, error) {
	op := newFindCommand(b, d)

	el := new(struct{ Value map[string]string })
	_, err := d.Client.ExecuteCmd(op, el)
	if err != nil {
		return nil, fmt.Errorf("error on find element by %+v\n error: %v", b, err)
	}

	eId := elementID(el.Value)

	return &Element{
		Id:       eId,
		Driver:   d,
		Selector: b,
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
