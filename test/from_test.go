package test

import (
	"testing"

	"github.com/mcsymiv/godriver/by"
	"github.com/mcsymiv/godriver/capabilities"
)

func TestFrom(t *testing.T) {
	d, tear := Driver(
		capabilities.HeadLess(),
		capabilities.Port("4444"),
	)

	defer tear()

	d.Open("https://google.com")
	nav := d.FindCss("[id='gb']")
	g := nav.From(by.Text("Gmail"))

	g.Click()
	d.Screenshot()

}
