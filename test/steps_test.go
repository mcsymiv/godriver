
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
  
		st.Cl("click on", "//*[@id='ceff842f-22cf-4f56-ac12-f30fa465761b']")
	
		st.Cl("click on", "//*[@id='tilesContainer']/div[1]/div[3]/div[2]")
	
		st.Cl("click on", "//*[@data-qa-id='format-shapes']")
	
		st.Cl("click on", "//*[@id='properties']/div[2]/div/div[1]/div[2]")
	
		st.Cl("click on", "//*[@id='properties']/div[2]/div/div[2]/div/div[1]/div[2]/div[1]/div/div/span")
	
		st.Cl("click on", "//*[@data-qa-id='insert-text']")
	
		st.Cl("click on", "//*[@id='Hyta_MCQHr-background']")
	
		st.Cl("click on", "//*[@id='root']/div/div/main/div[1]/div[1]/div[2]/div[2]/button")
	
		st.Cl("click on", "//*[@data-qa-id='back-to']")
	}