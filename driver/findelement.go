package driver

import (
	"fmt"

	"github.com/mcsymiv/godriver/by"
)

// F
func (d *Driver) F(selector string) *Element {
	w3cBy := by.Strategy(selector)

	el, err := f(w3cBy, d)
	if err != nil {
		panic(fmt.Errorf("unable to find element, got: %v\n", err))
	}

	return el
}

// TryF
// TODO:
// accepts slice of selectors and tryies to find element
// decrease timeout
// upd selector strategy if needed
// upd error handling
// add different find strategy
func (d *Driver) TryF(selectors []string) *Element {
	var el *Element

	for _, s := range selectors {
		w3cBy := by.Strategy(s)

		el, err := f(w3cBy, d)
		if err != nil {
			return nil
		}

		if el.Id != "" {
			return nil
		}

	}

	return el
}

func (d *Driver) Fs(selector string) []*Element {
	w3cBy := by.Css(selector)

	el, err := finds(w3cBy, d)
	if err != nil {
		return nil
	}

	return el
}
