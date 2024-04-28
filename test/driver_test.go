package test

import (
	// "fmt"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mcsymiv/godriver/config"
	"github.com/mcsymiv/godriver/driver"
)

func loginOkta(d *driver.Driver) {
	d.F("//*[@id='okta-signin-username']").Key(os.Getenv("OKTA_LOGIN"))
	d.F("//*[@id='okta-signin-password']").Key(os.Getenv("OKTA_PASS")).Key(driver.EnterKey)
}

func TestDeleteAccount(t *testing.T) {
	d, tear := Driver()
	defer tear()

	config.LoadEnv("../config", ".env")

	d.Url(os.Getenv("SUB_ENVIRONMENT_01"))
	loginOkta(d)

	d.F("//*[contains(@class, 'pagination_pageSize')]").Cl()
	d.F("200").Click()

	acc := "qa-dev01-135319"

	d.F(fmt.Sprintf("//*[text()='%s']/..//*[@data-qa-id='delete']", acc)).Cl()
	d.F("//*[text()='Confirm Delete']/../..//input").Key(acc)
	d.F("Yes").Cl()
}

func TestNewAccount(t *testing.T) {
	d, tear := Driver()
	defer tear()

	config.LoadEnv("../config", ".env")

	d.Url(os.Getenv("SUB_ENVIRONMENT_01"))
	loginOkta(d)

	acc := "qa-dev01-135319"

	d.F("Add Account").Cl()
	d.F("Customer Name *").Cl().Active().Key(acc)
	d.F("System Name *").Cl().Active().Key(acc)
	d.F("Sub Domain *").Cl().Active().Key(acc)
	// d.FindText("Built-in Authentication").Click()

	d.F("SMB").Cl()
	d.F("Enterprise").Cl()
	d.F("Create").Cl()
}

func TestDriver(t *testing.T) {
	d, tear := Driver()
	defer tear()

	repo := "/repository/download/"
	allure := ":id/allure-report.zip!/allure-report-test/index.html#suites"
	config.LoadEnv("../config", ".env")
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
	d.F("Log in using Azure Active Directory").Is().Cl()
	d.F("[id='i0116']").Key(os.Getenv("DOWNLOAD_LOGIN")).Key(driver.EnterKey)
	d.F("[id='i0118']").Key(os.Getenv("DOWNLOAD_PASS"))
	d.F("//input[@value='Увійти']").Is().Cl()
	// d.Find("//input[@value='Так']").IsDisplayed().Cl()
	d.F("Так").Is().Cl()
	d.F("Projects").Is().Cl()
	d.F("[id='search-projects']").Is().Key(testEnv)

	for _, sName := range sNames {
		// d.Find(fmt.Sprintf("//aside//span[contains(text(),'%s')]", sName)).IsDisplayed().Click()
		d.F("//*[@data-test='sidebar']").From(sName).Is().Cl()

		buildLinkRaw := d.F("(//*[@data-grid-root='true']//*[@data-test='ring-link'])[1]").Is().Attr("href")
		buildLink := strings.Join(strings.Split(buildLinkRaw, "/")[2:], "/")

		rLinks = append(rLinks, fmt.Sprintf("%s%s%s%s", host, repo, buildLink, allure))
	}

	for _, rLink := range rLinks {
		d.Url(rLink)
		time.Sleep(10 * time.Second)
		d.F("[data-tooltip='Download CSV']").Cl()
		time.Sleep(10 * time.Second)
	}
}
