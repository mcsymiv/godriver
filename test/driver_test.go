package test

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	"strings"
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

	repo := "/repository/download/"
	allure := ":id/allure-report.zip!/allure-report-test/index.html#suites"
	config.LoadEnv("../config", ".env")
	host := os.Getenv("DOWNLOAD_HOST")

	var rLinks []string
	sNames := []string{
		"Smoke (Concurrent tests)",
		"UI Regression (Concurrent tests)",
	}

	d.Open(fmt.Sprintf("%s%s", host, "/login.html"))
	d.FindX(".//a[text()='Log in using Azure Active Directory']").IsDisplayed().Click()
	d.FindCss("[id='i0116']").Key(os.Getenv("DOWNLOAD_LOGIN")).Key(driver.EnterKey)
	d.FindCss("[id='i0118']").Key(os.Getenv("DOWNLOAD_PASS"))
	d.FindX("//input[@value='Увійти']").IsDisplayed().Click()
	d.FindX("//input[@value='Так']").IsDisplayed().Click()
	d.FindX("//span[text()='Projects']").IsDisplayed().Click()
	d.FindCss("[id='search-projects']").IsDisplayed().Key("dev01")

	for _, sName := range sNames {
		d.FindX(fmt.Sprintf("//aside//span[contains(text(),'%s')]", sName)).IsDisplayed().Click()

		buildLinkRaw := d.FindX("(//*[@data-grid-root='true']//*[@data-test='ring-link'])[1]").IsDisplayed().Attribute("href")
		buildLink := strings.Join(strings.Split(buildLinkRaw, "/")[2:], "/")

		rLinks = append(rLinks, fmt.Sprintf("%s%s%s%s", host, repo, buildLink, allure))
	}

	for _, rLink := range rLinks {
		d.Open(rLink)
		time.Sleep(5 * time.Second)
		d.FindCss("[data-tooltip='Download CSV']").Click()
		time.Sleep(5 * time.Second)
	}
}
