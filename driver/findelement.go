package driver

import (
	"fmt"

	"github.com/mcsymiv/godriver/by"
)

// F
func (d *Driver) F(selector string) *Element {
	w3cBy := by.Strategy(selector)

	el, err := find(w3cBy, d)
	if err != nil {
		return &Element{
			ElementError: fmt.Errorf("unable to find element, got: %v", err),
		}
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

		el, err := find(w3cBy, d)
		if err != nil {
			return nil
		}

		if el.Id != "" {
			return el
		}

	}

	return el
}

func (d *Driver) Fs(selector string) *Element {
	w3cBy := by.Css(selector)

	el, err := find(w3cBy, d)
	if err != nil {
		return nil
	}

	return el
}

// Find
func (d *Driver) Find(selector string) *Element {
	w3cBy := by.Selector{
		Value: selector,
		Using: by.DefineStrategy(selector),
	}

	el, err := find(w3cBy, d)
	if err != nil {
		return nil
	}
	return el
}

// FindText
func (d *Driver) FindText(value string) *Element {
	w3cBy := by.Text(value)

	el, err := find(w3cBy, d)
	if err != nil {
		return nil
	}
	return el
}

// Find by w3c selector strategy
func (d *Driver) FindElement(by by.Selector) *Element {
	el, err := find(by, d)
	if err != nil {
		return nil
	}
	return el
}

func (d *Driver) Finds(selector string) []*Element {
	w3cBy := by.Selector{
		Value: selector,
		Using: by.DefineStrategy(selector),
	}

	el, err := finds(w3cBy, d)
	if err != nil {
		return nil
	}

	return el
}
