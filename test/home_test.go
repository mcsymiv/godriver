package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/godriver/capabilities"
)

func TestHome(t *testing.T) {
	d, tear := Driver(
		capabilities.Port("4445"),
	)
	defer tear()

	fmt.Println(os.Args)

	d.Url("http://192.168.0.1/")
	d.F("//*[@type='password']").Key(os.Getenv("HOME_PASS"))
	d.ExecuteScript("click.js", d.F("LOG IN").ElementIdentifier())
	time.Sleep(7 * time.Second)
	d.Cl("Clients")

	time.Sleep(7 * time.Second)
}
