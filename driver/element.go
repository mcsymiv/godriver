package driver

import (
	"fmt"
	"log"
	"net/http"
)

const (
	// LegacyWebElementIdentifier is the string constant used in the old Selenium 2 protocol
	// WebDriver JSON protocol that is the key for the map that contains an
	// unique element identifier.
	// This value is ignored in element id retreival
	LegacyWebElementIdentifier = "ELEMENT"

	// WebElementIdentifier is the string constant defined by the W3C Selenium 3 protocol
	// specification that is the key for the map that contains a unique element identifier.
	WebElementIdentifier = "element-6066-11e4-a52e-4f735466cecf"

	// ShadowRootIdentifier A shadow root is an abstraction used to identify a shadow root when
	// it is transported via the protocol, between remote and local ends.
	ShadowRootIdentifier = "shadow-6066-11e4-a52e-4f735466cecf"
)

const (
	ById              = "id" // not speciied by w3c
	ByXPath           = "xpath"
	ByLinkText        = "link text"
	ByPartialLinkText = "partial link text"
	ByName            = "name" // not specified by w3c
	ByTagName         = "tag name"
	ByClassName       = "class name" // not specified by w3c
	ByCssSelector     = "css selector"
)

// Empty
// Due to geckodriver bug: https://github.com/webdriverio/webdriverio/pull/3208
// "where Geckodriver requires POST requests to have a valid JSON body"
// Used in POST requests that don't require data to be passed by W3C
type Empty struct{}

// Element
// W3C WebElement
type Element struct {
	Id      string
	Client  *Client
	Session *Session
}

type JsonFindUsing struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

func (e Element) ElementIdentifier() map[string]string {
	return map[string]string{
		WebElementIdentifier: e.Id,
	}
}

func elementID(v map[string]string) string {
	id, ok := v[WebElementIdentifier]
	if !ok || id == "" {
		log.Println("Error on find element", v)
	}
	return id
}

func elementsID(v []map[string]string) []string {
	var els []string

	for _, el := range v {
		id, ok := el[WebElementIdentifier]
		if !ok || id == "" {
			log.Println("Error on find elements", v)
		}
		els = append(els, id)
	}

	return els
}

// Attribute
// Returns elements attribute value
func (e *Element) Attribute(a string) (string, error) {
	op := &Command{
		Path:   fmt.Sprintf("/element/%s/attribute/%s", e.Id, a),
		Method: http.MethodGet,
	}

	res, err := e.Client.ExecuteCommandStrategy(op)
	if err != nil {
		return "", err
	}

	attr := new(struct{ Value string })
	unmarshalData(res, attr)

	return attr.Value, nil
}