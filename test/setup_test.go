package test

import (
	"log"
	"os"

	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/config"
	"github.com/mcsymiv/godriver/driver"
	"github.com/mcsymiv/godriver/file"
	"github.com/xlzd/gotp"
)

func Driver(caps ...capabilities.CapabilitiesFunc) (driver.Driver, func()) {
	d := driver.NewDriver(caps...)
	if &d == nil {
		log.Fatal("Unable to start driver")
	}

	file.LoadEnv("../config", ".env")

	config.TestSetting = config.DefaultSetting()

	return d, func() {
		// teardown
		d.Quit()
		driver.OutFileLogs.Close()
		d.Service().Process.Kill()
	}
}

func loginOkta(d driver.Driver) {
	d.F("//*[@id='okta-signin-username']").Key(os.Getenv("OKTA_LOGIN"))
	d.F("//*[@id='okta-signin-password']").Key(os.Getenv("OKTA_PASS")).Key(driver.EnterKey)
	totp := gotp.NewDefaultTOTP(os.Getenv("OKTA_TOTP"))
	d.F("//*[@id='input59']").Key(totp.Now()).Key(driver.EnterKey)
}
