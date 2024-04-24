package test

import (
	"log"
	"testing"

	"github.com/mcsymiv/godriver/by"
	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/driver"
)

func fromDriver(caps ...capabilities.CapabilitiesFunc) (*driver.Driver, func()) {
	d := driver.NewDriver(caps...)
	if d == nil {
		log.Fatal("Unable to start driver")
	}
	return d, func() {
		// teardown
		d.Quit()
		driver.OutFileLogs.Close()
		d.Service().Process.Kill()
	}
}

func TestFrom(t *testing.T) {
	d, tear := fromDriver(
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
