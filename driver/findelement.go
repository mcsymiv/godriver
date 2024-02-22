package driver

import (
	"log"

	"github.com/mcsymiv/godriver/by"
)

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
	w3cBy := by.XPathTextStrategy(value)

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

func (d *Driver) FindX(selector string) *Element {
	by := by.Selector{
		Using: by.ByXPath,
		Value: selector,
	}

	el, err := find(by, d)
	if err != nil {
		return nil
	}
	return el
}

func (d *Driver) FindXs(selector string) []*Element {
	by := by.Selector{
		Using: by.ByXPath,
		Value: selector,
	}

	el, err := finds(by, d)
	if err != nil {
		return nil
	}
	return el
}

func (d *Driver) FindCss(selector string) *Element {
	by := by.Selector{
		Using: by.ByCssSelector,
		Value: selector,
	}

	el, err := find(by, d)
	if err != nil {
		log.Println(err)
		return nil
	}
	return el
}
