package steps

import (
	"fmt"
	"testing"

	"github.com/mcsymiv/godriver/driver"
)

type Test struct {
	*testing.T
	*driver.Driver
}

func (ts *Test) Url(arg string) {
	var name string = fmt.Sprintf("open %s", arg)
	ts.T.Run(name, func(t *testing.T) {
		d := ts.Driver.Url(arg)
		if d == nil {
			t.Fatal("unable to open url")
		}
	})
}

func (ts *Test) Cl(arg string) *driver.Element {
	var name string = fmt.Sprintf("click on %s", arg)
	var el *driver.Element

	ts.T.Run(name, func(t *testing.T) {
		el = ts.Driver.F(arg)
		if el == nil {
			t.Fatal("unable to find element")
		}

		el = el.Is()
		if el == nil {
			t.Fatal("element not displayed")
		}

		el = el.Click()
		if el == nil {
			t.Fatal("unable to click on element")
		}
	})

	return el
}

func (ts *Test) Is(name string, arg string) {
	ts.T.Run(name, func(t *testing.T) {
		el := ts.Driver.F(arg)
		if el == nil {
			t.Fatal("unable to find element")
		}

		el = el.Is()
		if el == nil {
			t.Fatal("element not displayed")
		}
	})
}
