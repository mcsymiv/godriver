package test

import (
	"testing"

	"github.com/mcsymiv/godriver/capabilities"
)

func TestJs(t *testing.T) {
	d, tear := Driver(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	d.Url("http://google.com")
	el := d.F("[aria-label='Google apps']")
	d.ExecuteScript("click.js", el.ElementIdentifier())
}
