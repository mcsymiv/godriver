package driver

import (
	"fmt"
	"net/http"

	"github.com/mcsymiv/godriver/by"
)

const (
	fromElementPath  = "/element/%s/element"
	fromElementsPath = "/element/%s/elements"
)

func from(by by.Selector, e *Element) (*Element, error) {
	op := &Command{
		Path:           fromElementPath,
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategies: []CommandExecutor{newFindStrategy(e.Driver)},
	}

	el := new(struct{ Value map[string]string })
	_, err := e.Driver.Client.ExecuteCmd(op, el)
	if err != nil {
		return nil, fmt.Errorf("error on find element from: %v", err)
	}

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

// FromText
// Finds Element from receiver Element by.Text
func (e *Element) FromText(txt string) *Element {
	el, err := from(by.Text(txt), e)
	if err != nil {
		return nil
	}

	return el
}

func (e *Element) Froms(by by.Selector) []*Element {
	op := &Command{
		Path:           fromElementsPath,
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodPost,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
		Strategies: []CommandExecutor{newFindStrategy(e.Driver)},
	}

	el := new(struct{ Value []map[string]string })
	e.Driver.Client.ExecuteCmd(op, el)
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
