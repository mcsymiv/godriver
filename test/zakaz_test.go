package test

import (
	"fmt"
	"testing"

	"github.com/mcsymiv/godriver/by"
	"github.com/mcsymiv/godriver/capabilities"
)

func TestZakaz(t *testing.T) {
	d, tear := Driver(
		capabilities.HeadLess(),
		capabilities.Port("4444"),
	)
	defer tear()

	t.Run("zakaz", func(t *testing.T) {
		d.Url("https://zakaz.ua/en/")
		d.F("[data-marker='NOVUS']").Cl()
		d.Tab(1).F("[data-marker='Close popup']").Cl()
		d.F("Grocery").Cl()
		d.F("//h1[text()='Grocery']").Is()

		els := d.F("[id='PageWrapBody_desktopMode']").Froms(by.Css("[data-testid='product_tile_inner']"))

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
