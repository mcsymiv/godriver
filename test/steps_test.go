package test

import (
	"testing"

	"github.com/mcsymiv/godriver/steps"
)

func TestGeneratedSteps(t *testing.T) {
	d, tear := Driver()
	defer tear()

	st := steps.Test{t, d}

	st.Url("open page", "https://google.com")

	st.Cl("Search").Key("QAAuto Asset Import02")

	st.Cl("QAAuto Asset Import02")

	st.Cl("//*[@data-qa-id='format-shapes']")

	st.Cl("Permissions")

	st.Cl("Manage Pages")

	st.Cl("//*[@data-qa-id='insert-text']")

	st.Cl("//*[@id='Hyta_MCQHr-background']")

	st.Cl("Save")

	st.Cl("//*[@data-qa-id='back-to']")
}
