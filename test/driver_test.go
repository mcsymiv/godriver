package test

import (
	"log"
	"testing"
	"time"

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
		d.Service().Process.Kill()
	}
}

func TestDriver(t *testing.T) {
	d, tear := Driver()
	defer tear()

	d.Open("https://google.com")
	d.FindX("//*[@id='APjFqb']").Click()
	d.Find(driver.By{Using: driver.ByXPath, Value: "//*[@id='APjFqb']"}).Click()
	time.Sleep(5 * time.Second)
}
