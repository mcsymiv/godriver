package test

import (
	"log"

	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/driver"
)

func Driver(caps ...capabilities.CapabilitiesFunc) (*driver.Driver, func()) {
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
