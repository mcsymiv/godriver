package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type AutoGenerated struct {
	Title string `json:"title"`
	Steps []struct {
		Type              string `json:"type"`
		Width             int    `json:"width,omitempty"`
		Height            int    `json:"height,omitempty"`
		DeviceScaleFactor int    `json:"deviceScaleFactor,omitempty"`
		IsMobile          bool   `json:"isMobile,omitempty"`
		HasTouch          bool   `json:"hasTouch,omitempty"`
		IsLandscape       bool   `json:"isLandscape,omitempty"`
		URL               string `json:"url,omitempty"`
		AssertedEvents    []struct {
			Type  string `json:"type"`
			URL   string `json:"url"`
			Title string `json:"title"`
		} `json:"assertedEvents,omitempty"`
		Target    string     `json:"target,omitempty"`
		Selectors [][]string `json:"selectors,omitempty"`
		OffsetY   float32    `json:"offsetY,omitempty"`
		OffsetX   float32    `json:"offsetX,omitempty"`
	} `json:"steps"`
}

func CreateTest() {

	file, err := os.Open("../artifacts/records/rec_2.json")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	byteValue, _ := io.ReadAll(file)

	var at AutoGenerated

	err = json.Unmarshal(byteValue, &at)
	if err != nil {
		log.Fatal("unable to unmarshal record.json", err)
	}

	gF, err := os.OpenFile("../test/demo_test.go", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer gF.Close()

	var testStr string = `
	package test

  import "testing"

  func Test%s(t *testing.T) {
  	d, tear := Driver()
  	defer tear()

		d.Url("https://google.com")

  ` // operand expected, found }

	gF.WriteString(fmt.Sprintf(testStr, "Generated"))

	var clickStr string = `
		d.F("%s").Click()
	`

	for _, st := range at.Steps {
		if st.Type == "click" && len(st.Selectors) > 0 {
			for _, s := range st.Selectors {
				if strings.Contains(s[0], "xpath") {
					xFormated := strings.ReplaceAll(s[0], "\"", "'")
					xFormated = strings.ReplaceAll(xFormated, "xpath/", "")
					gF.WriteString(fmt.Sprintf(clickStr, xFormated))
				}
			}
		}
	}

	gF.WriteString("}")
}