package test

import (
	"log"
	"testing"

	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/driver"
)

func Driver(caps ...capabilities.CapabilitiesFunc) (*driver.Driver, func()) {
	d := driver.NewDriver(caps...)
	if d == nil {
		log.Fatal("Unable to start driver")
	}
	// tear down later
	return d, func() {
		// tear-down code here
		d.Quit()
		d.Service().Process.Kill()
	}
}

func TestDriver(t *testing.T) {
	d, tear := Driver()
	defer tear()

	d.Open("https://google.com")
	el, _ := d.FindElement("[id='APjFqb']")
	el.Click()
}
