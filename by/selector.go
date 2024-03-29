package by

import (
	"fmt"
	"strings"
)

// w3c Locator strategies
// src: https://www.w3.org/TR/webdriver2/#locator-strategies
const (
	ByXPath           = "xpath"
	ByLinkText        = "link text"
	ByPartialLinkText = "partial link text"
	ByTagName         = "tag name"
	ByCssSelector     = "css selector"
)

type Selector struct {
	Using, Value string
}

func DefineStrategy(s string) string {
	if strings.Contains(s, "/") {
		return ByXPath
	}

	return ByCssSelector
}

// XPathTextStrategy
// text/value based find strategy
// TODO: remove duplicates in fmt.Sprinf
func XPathTextStrategy(value string) Selector {
	return Selector{
		Using: ByXPath,
		Value: fmt.Sprintf("//*[text()='%s'] | //*[@placeholder='%s'] | //*[@value='%s']", value, value, value),
	}
}
