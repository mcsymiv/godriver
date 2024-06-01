package driver

import (
	"net/http"

	"github.com/mcsymiv/godriver/by"
)

func f(b by.Selector, d Driver) (Element, error) {
	el := new(struct{ Value map[string]string })

	st := findStrategyV2{
		Driver: d,
		Command: Command{
			Path:         PathElementFind,
			Method:       http.MethodPost,
			ResponseData: el,
			Data: marshalData(&JsonFindUsing{
				Using: b.Using,
				Value: b.Value,
			}),
		},
	}

	st.execute()

	eId := elementID(el.Value)

	return Element{
		Id:       eId,
		Driver:   d,
		Selector: b,
	}, nil
}

func finds(b by.Selector, d Driver) ([]Element, error) {
	el := new(struct{ Value []map[string]string })

	st := findStrategyV2{
		Driver: d,
		Command: Command{
			Path:         PathElementsFind,
			Method:       http.MethodPost,
			ResponseData: el,
			Data: marshalData(&JsonFindUsing{
				Using: b.Using,
				Value: b.Value,
			}),
		},
	}

	st.execute()

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
