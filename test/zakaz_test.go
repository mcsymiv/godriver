package test

import (
	"fmt"
	"testing"

	"github.com/mcsymiv/godriver/by"
	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/driver"
)

type Test struct {
	*testing.T
	*driver.Driver
}

func (ts *Test) url(s string, arg string) {
	ts.T.Run(s, func(t *testing.T) {
		ts.Driver.Open(arg)
	})
}

func (ts *Test) click(s string, arg string) {
	ts.T.Run(s, func(t *testing.T) {
		el := ts.Driver.Find(arg)
		el.Click()
	})
}

func TestZakaz(t *testing.T) {
	d, tear := Driver(
		capabilities.HeadLess(),
		capabilities.Port("4444"),
	)
	defer tear()

	test := &Test{t, d}

	test.url("open url", "https://zakaz.ua/en/")
	test.click("novus banner", "[data-marker='NOVUS']")

	t.Run("zakaz", func(t *testing.T) {
		d.SwitchToTab(1)
		d.Find("[data-marker='Close popup']").Click()

		d.FindText("Grocery").Click()
		d.Find("//h1[text()='Grocery']").IsDisplayed()

		els := d.Find("[id='PageWrapBody_desktopMode']").Froms(by.Css("[data-testid='product_tile_inner']"))

		var products [][]string

		for _, el := range els {
			var product []string = []string{}

			pn := el.From("[data-testid='product_tile_title']")

			product = append(product, pn.Text())
			products = append(products, product)
		}

		fmt.Println(products)
	})

}
