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

func checkSubstrings(str string, subs ...string) (bool, int) {
	matches := 0
	isCompleteMatch := true

	for _, sub := range subs {
		if strings.Contains(str, sub) {
			matches += 1
		} else {
			isCompleteMatch = false
		}
	}

	return isCompleteMatch, matches
}

func Strategy(value string) Selector {
	if strings.Contains(value, "/") {
		return Selector{
			Value: value,
			Using: ByXPath,
		}
	}

	// if ok, m := checkSubstrings(value, ".", "#", "[", "]"); ok || m > 0 {
	if strings.Contains(value, "[") {
		return Selector{
			Value: value,
			Using: ByCssSelector,
		}
	}

	return xPathTextStrategy(value)
}

// XPathTextStrategy
// text/value based find strategy
func xPathTextStrategy(value string) Selector {
	return Selector{
		Using: ByXPath,
		Value: fmt.Sprintf("//*[text()='%[1]s'] | //*[@placeholder='%[1]s'] | //*[@value='%[1]s']", value),
	}
}

func Text(value string) Selector {
	return xPathTextStrategy(value)
}

func Css(value string) Selector {
	return Selector{
		Using: ByCssSelector,
		Value: value,
	}
}
