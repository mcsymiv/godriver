package driver

import (
	"github.com/mcsymiv/godriver/by"
)

// F
func (d *Driver) F(selector string) *Element {
	w3cBy := by.Strategy(selector)

	el, err := find(w3cBy, d)
	if err != nil {
		return nil
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
