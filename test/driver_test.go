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
	testEnv := "dev01"

	var rLinks []string
	sNames := []string{
		// os.Getenv("SUITE_NAME_1"),
		os.Getenv("SUITE_NAME_2"),
		os.Getenv("SUITE_NAME_3"),
		os.Getenv("SUITE_NAME_4"),
		// os.Getenv("SUITE_NAME_5"),
		// os.Getenv("SUITE_NAME_6"),
		os.Getenv("SUITE_NAME_7"),
		os.Getenv("SUITE_NAME_8"),
		os.Getenv("SUITE_NAME_9"),
		// os.Getenv("SUITE_NAME_10"),
	}

	d.Open(fmt.Sprintf("%s%s", host, "/login.html"))
	d.Find(".//a[text()='Log in using Azure Active Directory']").IsDisplayed().Click()
	d.Find("[id='i0116']").Key(os.Getenv("DOWNLOAD_LOGIN")).Key(driver.EnterKey)
	d.Find("[id='i0118']").Key(os.Getenv("DOWNLOAD_PASS"))
	d.Find("//input[@value='Увійти']").IsDisplayed().Click()
	d.Find("//input[@value='Так']").IsDisplayed().Click()
	d.Find("//span[text()='Projects']").IsDisplayed().Click()
	d.Find("[id='search-projects']").IsDisplayed().Key(testEnv)

	for _, sName := range sNames {
		d.Find(fmt.Sprintf("//aside//span[contains(text(),'%s')]", sName)).IsDisplayed().Click()

		buildLinkRaw := d.Find("(//*[@data-grid-root='true']//*[@data-test='ring-link'])[1]").IsDisplayed().Attribute("href")
		buildLink := strings.Join(strings.Split(buildLinkRaw, "/")[2:], "/")

		rLinks = append(rLinks, fmt.Sprintf("%s%s%s%s", host, repo, buildLink, allure))
	}

	for _, rLink := range rLinks {
		d.Open(rLink)
		time.Sleep(5 * time.Second)
		d.Find("[data-tooltip='Download CSV']").Click()
		time.Sleep(5 * time.Second)
	}
}

func TestStep(t *testing.T) {
	d, tear := Driver(
		capabilities.ImplicitWait(10000),
		capabilities.PageLoadStrategy("eager"),
	)
	defer tear()

	config.LoadEnv("../config", ".env")

	d.Open(os.Getenv("LOCAL_HOST"))
	d.Find("//*[@id='btn-sign-in']").Key(driver.EnterKey)

	d.FindText("My Account").Click()
	d.FindText("Manage Asset Types").Click()
	d.FindText("Add Root Item").Click()
	d.FindText("Name *").Click()
	d.GetActive().Key("worked")

	time.Sleep(4 * time.Second)
}
