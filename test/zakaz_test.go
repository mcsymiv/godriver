package test

import (
	"testing"

	"github.com/mcsymiv/godriver/capabilities"
)

func TestZakaz(t *testing.T) {
	d, tear := Driver(
		capabilities.HeadLess(),
		capabilities.Port("4444"),
	)

	defer tear()

	d.Open("https://zakaz.ua/en/")
	d.FindCss("[data-marker='NOVUS']").IsDisplayed().Click()
}
