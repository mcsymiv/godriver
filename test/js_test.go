package test

import (
	"testing"
	"time"

	"github.com/mcsymiv/godriver/capabilities"
)

func TestJs(t *testing.T) {
	d, tear := Driver(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	d.Url("http://google.com")
	time.Sleep(5 * time.Second)
	el := d.F("[aria-label='Google apps']")
	// d.Script("arguments[0].click();", el.ElementIdentifier())
	d.ExecuteScript("click.js", el.ElementIdentifier())
	time.Sleep(5 * time.Second)
}
