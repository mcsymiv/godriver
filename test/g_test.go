package test

import (
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/driver"
)

func TestG(t *testing.T) {
	d, tear := Driver(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	d.Url("https://console.cloud.google.com/")
	d.F("//*[@id='identifierId']").Key(os.Getenv("G_USER")).Key(driver.EnterKey)
	d.ClickJs("//*[contains(text(), 'Show password')]")
	d.SetValueJs("//*[@id='password']//input", os.Getenv("G_PASS"))
	time.Sleep(1 * time.Second)
	d.ClickJs("//*[text()='Next']")
	d.Cl("APIs and services")
	d.Cl(" Credentials ")
	d.Cl("//*[text()=' OAuth 2.0 Client IDs ']/../..//*[@data-mat-icon-name='delete']")
	d.Cl("Delete")
	d.SetValueJs("//mat-form-field//input", "DELETE")
	d.Cl("Delete")
	d.ClickJs("//*[text()='\n  Create credentials']")
	d.Cl(" OAuth client ID ")
	d.Cl("Application type")
	d.Cl("Desktop app")
	d.Cl("//button//span[contains(text(),'Create')]")
	d.Cl("Download JSON")

	time.Sleep(4 * time.Second)
}

func TestG2(t *testing.T) {
	d, tear := Driver(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	d.Url("https://console.cloud.google.com/")
	d.F("//*[@id='identifierId']").Key(os.Getenv("G_USER")).Key(driver.EnterKey)
	d.ClickJs("//*[contains(text(), 'Show password')]")
	d.SetValueJs("//*[@id='password']//input", os.Getenv("G_PASS"))
	time.Sleep(1 * time.Second)
	d.ClickJs("//*[text()='Next']")
	d.Cl("APIs and services")
	d.Cl(" Credentials ")
	time.Sleep(1 * time.Second)
	d.ClickJs("//*[text()='\n  Create credentials']")
	d.Cl(" OAuth client ID ")
	d.Cl("Application type")
	d.Cl("Desktop app")
	d.Cl("//button//span[contains(text(),'Create')]")
	d.Cl("Download JSON")

	time.Sleep(4 * time.Second)
}
