package driver

import (
	"net/http"

	"github.com/mcsymiv/godriver/by"
)

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
		Strategy: &findStrategy{d},
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
