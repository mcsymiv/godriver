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
		Path:   PathElementFind,
		Method: http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),

		Strategy: newFindStrategy(d),
	}
}

func find(b by.Selector, d *Driver) (*Element, error) {
	op := newFindCommand(b, d)

	el := new(struct{ Value map[string]string })
	d.Client.ExecuteCommand(op, el)

	eId := elementID(el.Value)

	return &Element{
		Id:       eId,
		Driver:   d,
		Selector: b,
	}, nil
}

func f(b by.Selector, d *Driver) (*Element, error) {
	el := new(struct{ Value map[string]string })
	d.Client.ExecuteCommand(&Command{
		Path:   PathElementFind,
		Method: http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: b.Using,
			Value: b.Value,
		}),

		Strategy: &findStrategy{
			driver: d,
		},
	}, el)

	eId := elementID(el.Value)

	return &Element{
		Id:       eId,
		Driver:   d,
		Selector: b,
	}, nil
}

func finds(by by.Selector, d *Driver) ([]*Element, error) {
	op := &Command{
		Path:   PathElementsFind,
		Method: http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategy: newFindStrategy(d),
	}

	el := new(struct{ Value []map[string]string })
	d.Client.ExecuteCommand(op, el)
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
