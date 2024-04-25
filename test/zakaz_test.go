package test

import (
	"fmt"
	"testing"

	"github.com/mcsymiv/godriver/by"
	"github.com/mcsymiv/godriver/capabilities"
)

func TestZakaz(t *testing.T) {
	d, tear := Driver(
		capabilities.Port("4444"),
	)

	defer tear()

	d.Open("https://zakaz.ua/en/")
	d.FindCss("[data-marker='NOVUS']").Click()
	d.SwitchToTab(1)
	d.Find("[data-marker='Close popup']").Click()

	d.FindText("Grocery").Click()
	d.FindX("//h1[text()='Grocery']").IsDisplayed()

	els := d.Find("[id='PageWrapBody_desktopMode']").Froms(by.Css("[data-testid='product_tile_inner']"))

	var products [][]string

	for _, el := range els {
		var product []string = []string{}

		pn := el.From(by.Css("[data-testid='product_tile_title']"))

		product = append(product, pn.Text())
		products = append(products, product)
	}

	fmt.Println(products)
}
