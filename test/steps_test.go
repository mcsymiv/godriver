package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/steps"
)

func TestGeneratedSteps(t *testing.T) {
	d, tear := Driver(
		capabilities.HeadLess(),
	)
	defer tear()

	st := steps.Test{t, d}

	fmt.Println(os.Getenv("SUB_ENVIRONMENT"))

	st.Url("https://google.com")

	// st.Cl("//*[@id='ceff842f-22cf-4f56-ac12-f30fa465761b']")
	//
	// st.Cl("QAAuto Asset Import02")
	//
	// st.Cl("//*[@data-qa-id='format-shapes']")
	//
	// st.Cl("Permissions")
	//
	// st.Cl("Manage Pages")
	//
	// st.Cl("//*[@data-qa-id='insert-text']")
	//
	// st.Cl("//*[@id='Hyta_MCQHr-background']")
	//
	// st.Cl("Save")
	//
	// st.Cl("//*[@data-qa-id='back-to']")
}
