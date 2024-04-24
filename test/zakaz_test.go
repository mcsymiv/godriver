package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mcsymiv/godriver/by"
	"github.com/mcsymiv/godriver/capabilities"
)

func TestZakaz(t *testing.T) {
	d, tear := Driver(
		capabilities.HeadLess(),
		capabilities.Port("4444"),
	)

	defer tear()

	d.Open("https://zakaz.ua/en/")
	d.FindCss("[data-marker='NOVUS']").Click()
	d.SwitchToTab(1)
	d.FindText("Grocery").Click()
	time.Sleep(2 * time.Second)
	els := d.Find("[id='PageWrapBody_desktopMode']").Froms(by.Selector{
		Using: by.ByCssSelector,
		Value: "[data-testid='product_tile_inner']",
	})

	for _, el := range els {
		oldP := el.From(by.Selector{
			Using: by.ByCssSelector,
			Value: "[data-marker='Price']",
		})

		fmt.Println(oldP.Text())
	}

	d.Screenshot()
}
