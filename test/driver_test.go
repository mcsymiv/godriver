package test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/driver"
)

func TestDriver(t *testing.T) {
	d, tear := Driver(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	repo := "/repository/download/"
	allure := ":id/allure-report.zip!/allure-report-test/index.html#suites"
	host := os.Getenv("DOWNLOAD_HOST")
	testEnv := "dev01"

	var rLinks []string
	sNames := []string{
		// os.Getenv("SUITE_NAME_1"), // smoke
		os.Getenv("SUITE_NAME_2"), // regress
		os.Getenv("SUITE_NAME_3"), // single
		os.Getenv("SUITE_NAME_4"), // m
		// os.Getenv("SUITE_NAME_5"), // ol
		// os.Getenv("SUITE_NAME_6"), // hil
		// os.Getenv("SUITE_NAME_7"), // gm
		// os.Getenv("SUITE_NAME_8"), // business
		// os.Getenv("SUITE_NAME_9"), // visual
		// os.Getenv("SUITE_NAME_10"), 		// iframe
	}

	d.Url(fmt.Sprintf("%s%s", host, "/login.html"))
	d.F("Log in using Azure Active Directory").Click()
	d.F("[aria-label^='Ending with']").Key(os.Getenv("DOWNLOAD_LOGIN")).Key(driver.EnterKey)
	d.F("[aria-label^='Enter the password']").Key(os.Getenv("DOWNLOAD_PASS")).Key(driver.EnterKey)
	d.F("Yes").Click()
	d.F("Projects").Click()
	d.F("[id='search-projects']").Key(testEnv)

	for _, sName := range sNames {
		log.Println(sName)
		d.F(fmt.Sprintf("//*[@data-test='sidebar']//span[contains(text(),'%s')]", sName)).Is().Click()
		d.F(fmt.Sprintf("//h1//*[text()='%s']", sName)).Is()

		buildLinkRaw := d.F("(//*[@data-grid-root='true']//*[@data-test='ring-link'])[1]").Attr("href")
		buildLink := strings.Join(strings.Split(buildLinkRaw, "/")[2:], "/")

		log.Println(buildLink)
		rLinks = append(rLinks, fmt.Sprintf("%s%s%s%s", host, repo, buildLink, allure))
	}

	for _, rLink := range rLinks {
		d.Url(rLink)
		time.Sleep(10 * time.Second)
		d.F("[data-tooltip='Download CSV']").Click()
		time.Sleep(10 * time.Second)
	}
}
