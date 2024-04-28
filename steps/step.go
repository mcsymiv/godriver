package steps

import (
	"testing"

	"github.com/mcsymiv/godriver/driver"
)

type Test struct {
	*testing.T
	*driver.Driver
}

func (ts *Test) Url(s string, arg string) {
	ts.T.Run(s, func(t *testing.T) {
		d := ts.Driver.Url(arg)
		if d == nil {
			t.Fatal("unable to open url")
		}
	})
}

func (ts *Test) Cl(s string, arg string) {
	ts.T.Run(s, func(t *testing.T) {
		el := ts.Driver.F(arg)
		if el == nil {
			t.Fatal("unable to find element")
		}

		el = el.IsDisplayed()
		if el == nil {
			t.Fatal("element not displayed")
		}

		el = el.Cl()
		if el == nil {
			t.Fatal("unable to click on element")
		}
	})
}

func (ts *Test) Is(name string, arg string) {
	ts.T.Run(name, func(t *testing.T) {
		el := ts.Driver.F(arg)
		if el == nil {
			t.Fatal("unable to find element")
		}

		el = el.IsDisplayed()
		if el == nil {
			t.Fatal("element not displayed")
		}
	})
}
