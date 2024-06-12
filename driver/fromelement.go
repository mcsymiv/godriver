package driver

import (
	"net/http"

	"github.com/mcsymiv/godriver/by"
)

func from(b by.Selector, e Element) (Element, error) {
	el := new(struct{ Value map[string]string })

	e.Driver.execute(findStrategyV2{Command{
		Path:           PathElementFromElement,
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodPost,
		ResponseData:   el,
		Data: marshalData(&JsonFindUsing{
			Using: b.Using,
			Value: b.Value,
		}),
	}})

	eId := elementID(el.Value)

	return Element{
		Id:       eId,
		Driver:   e.Driver,
		Selector: b,
	}, nil

}

// From
// Finds Element from receiver Element
func (e Element) From(s string) Element {
	by := by.Strategy(s)

	el, err := from(by, e)
	if err != nil {
		return Element{}
	}

	return el
}

func (e Element) Froms(by by.Selector) []Element {
	el := new(struct{ Value []map[string]string })

	e.Driver.execute(findStrategyV2{Command{
		Path:           PathElementsFromElement,
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodPost,
		ResponseData:   el,
		Data: marshalData(&JsonFindUsing{
			Using: by.Using,
			Value: by.Value,
		}),
	}})

	elementsIds := elementsID(el.Value)

	var els []Element

	for _, id := range elementsIds {
		els = append(els, Element{
			Id:     id,
			Driver: e.Driver,
		})
	}

	return els
}
