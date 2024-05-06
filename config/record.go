package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mcsymiv/go-brand/file"
)

// AutoGenerated
// Chrome struct for record
type AutoGenerated struct {
	Title              string               `json:"title"`
	AutoGeneratedSteps []AutoGeneratedSteps `json:"steps"`
}

// AutoGeneratedSteps
// Chrome struct of Steps record
type AutoGeneratedSteps struct {
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
}

// generatedStepSelectors
type GeneratedStepSelectors struct {
	step, css, text, xpath string
}

func readRecord() *AutoGenerated {
	f, err := file.Find("../artifacts", "record_3.json")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(f.FilePath)
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

	return &at
}

// convertSelectors
// formats chrome steps array selectors type to struct
func convertSelectors(at []AutoGeneratedSteps) []*GeneratedStepSelectors {
	var genSelectors []*GeneratedStepSelectors = []*GeneratedStepSelectors{}

	for _, st := range at {
		if st.Type == "click" && len(st.Selectors) > 0 {
			var genSelector *GeneratedStepSelectors = &GeneratedStepSelectors{
				step: st.Type,
			}

			for _, s := range st.Selectors {
				if strings.Contains(s[0], "xpath/") {
					xFormated := strings.ReplaceAll(s[0], "\"", "'")
					xFormated = strings.ReplaceAll(xFormated, "xpath/", "")
					genSelector.xpath = xFormated
				}

				if strings.Contains(s[0], "text/") {
					tFormated := strings.ReplaceAll(s[0], "text/", "")
					genSelector.text = tFormated
				}

				// ignore chrome selectors type if specified in record
				if strings.Contains(s[0], "aria") || strings.Contains(s[0], "pierce") {
					continue
				}

				genSelector.css = s[0]
			}

			genSelectors = append(genSelectors, genSelector)
		}
	}

	return genSelectors
}

func CreateTest() {

	at := readRecord()

	// TODO: find a way to remove relative path
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
	genS := convertSelectors(at.AutoGeneratedSteps)

	for _, st := range genS {
		gF.WriteString(fmt.Sprintf(clickStr, "click on", st.xpath))
	}

	gF.WriteString("}")
}

func CreateSteps() {
	at := readRecord()

	// TODO: find a way to remove relative path
	gF, err := os.OpenFile("../test/steps_test.go", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer gF.Close()

	var testStr string = `
	package test

  import (
		"testing"

		"github.com/mcsymiv/godriver/steps"
	)

  func Test%s(t *testing.T) {
  	d, tear := Driver()
  	defer tear()

		st := steps.Test{t, d}

		st.Url("open page", "https://google.com")
  ` // operand expected, found },

	gF.WriteString(fmt.Sprintf(testStr, "GeneratedSteps"))

	var clickStr string = `
		st.Cl("%s", "%s")
	` // step.Click method with %description, and %selector

	genS := convertSelectors(at.AutoGeneratedSteps)

	for _, st := range genS {
		gF.WriteString(fmt.Sprintf(clickStr, "click on", st.xpath))
	}

	gF.WriteString("}") // close Test%Name brakets
}
