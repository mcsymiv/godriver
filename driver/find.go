package driver

import (
	"net/http"

	"github.com/mcsymiv/godriver/by"
)

func f(b by.Selector, d Driver) (Element, error) {
	el := new(struct{ Value map[string]string })

	d.execute(retryStrategy{Command{
		Path:         PathElementFind,
		Method:       http.MethodPost,
		ResponseData: el,
		Data: marshalData(&JsonFindUsing{
			Using: b.Using,
			Value: b.Value,
		}),
	}})

	eId := elementID(el.Value)

	return Element{
		Id:       eId,
		Driver:   d,
		Selector: b,
	}, nil
}

func finds(b by.Selector, d Driver) ([]Element, error) {
	el := new(struct{ Value []map[string]string })

	d.execute(retryStrategy{Command{
		Path:         PathElementsFind,
		Method:       http.MethodPost,
		ResponseData: el,
		Data: marshalData(&JsonFindUsing{
			Using: b.Using,
			Value: b.Value,
		}),
	}})

	elementsId := elementsID(el.Value)

	var els []Element

	for _, id := range elementsId {
		els = append(els, Element{
			Id:     id,
			Driver: d,
		})
	}

	return els, nil
}
