package driver

import (
	"net/http"

	"github.com/mcsymiv/godriver/by"
)

func from(by by.Selector, e *Element) (*Element, error) {
	op := &Command{
		Path:           PathElementFromElement,
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategy: &findStrategy{e.Driver},
	}

	el := new(struct{ Value map[string]string })
	e.Driver.Client.ExecuteCommand(op, el)

	eId := elementID(el.Value)

	return &Element{
		Id:       eId,
		Driver:   e.Driver,
		Selector: by,
	}, nil

}

// From
// Finds Element from receiver Element
func (e *Element) From(s string) *Element {
	by := by.Strategy(s)

	el, err := from(by, e)
	if err != nil {
		return nil
	}

	return el
}

func (e *Element) Froms(by by.Selector) []*Element {
	op := &Command{
		Path:           PathElementsFromElement,
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategy: &findStrategy{e.Driver},
	}

	el := new(struct{ Value []map[string]string })
	e.Driver.Client.ExecuteCommand(op, el)
	elementsIds := elementsID(el.Value)

	var els []*Element

	for _, id := range elementsIds {
		els = append(els, &Element{
			Id:     id,
			Driver: e.Driver,
		})
	}

	return els
}
