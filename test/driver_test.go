package test

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/config"
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
	d, tear := Driver(
		capabilities.ImplicitWait(10000),
		capabilities.PageLoadStrategy("eager"),
	)
	defer tear()

	config.LoadEnv("../config", ".env")
	host := os.Getenv("DOWNLOAD_HOST")
	d.Open(fmt.Sprintf("%s%s", host, "/login.html"))

	d.FindX(".//a[text()='Log in using Azure Active Directory']").IsDisplayed().Click()
	d.FindCss("[id='i0116']").Key(os.Getenv("DOWNLOAD_LOGIN")).Key(driver.EnterKey)
	d.FindCss("[id='i0118']").Key(os.Getenv("DOWNLOAD_PASS"))

	time.Sleep(1 * time.Second)
	d.FindCss("[id='idSIButton9']").IsDisplayed().Click()
	d.FindX("//*[@value='Так']").IsDisplayed().Click()
	d.FindX(".//span[text()='Projects']").Click()

	time.Sleep(2 * time.Second)
	d.FindCss("[id='search-projects']").Key("dev01")
	d.FindX(".//aside//span[contains(text(),'UI Regression (Concurrent tests)')]").Click()
	els := d.FindXs("//*[@data-grid-root='true']//*[@data-test='ring-link']")
	els[0].Click()

	time.Sleep(10 * time.Second)
	d.FindX("(.//*[text()='Allure Report'])[1]").Click()

	time.Sleep(20 * time.Second)
	d.FindCss("[id*='iFrameResizer']").SwitchFrame()
	d.FindCss("[id='iframe']").SwitchFrame()

	time.Sleep(1 * time.Second)
	d.FindX(".//ul[@class='side-nav__menu']//div[text()='Suites']").Click()
	d.FindCss("[data-tooltip='Download CSV']").Click()

	time.Sleep(5 * time.Second)
}
