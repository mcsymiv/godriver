package driver

import (
	"log"
	"net/http"
	"os"

	"github.com/mcsymiv/godriver/config"
	"github.com/mcsymiv/godriver/file"
)

func readFile(name string) ([]byte, error) {
	c, err := os.ReadFile(name)
	if err != nil {
		log.Println("error on read file", err)
		return nil, err
	}

	return c, nil
}

// ClickJs
// Combines selenium selector strategy
// And Find element method with JS click
func (d *Driver) ClickJs(selector string) {
	el := d.F(selector)

	args := []interface{}{el.ElementIdentifier()}

	f := file.FindFile(config.TestSetting.JsFilesPath, "click.js")

	c, err := readFile(f)
	if err != nil {
		log.Println("error on file read in click.js", err)
		return
	}

	_, err = executeScriptSync(d, string(c), args)
	if err != nil {
		log.Println("error on execute script", err)
		return
	}

	return
}

// SetValueJs 
// Combines selenium selector strategy
// And Find element method with JS set value
func (d *Driver) SetValueJs(selector, value string) {
	el := d.F(selector)

	args := []interface{}{el.ElementIdentifier(), value}

	f := file.FindFile(config.TestSetting.JsFilesPath, "setValue.js")

	c, err := readFile(f)
	if err != nil {
		log.Println("error on file read in click.js", err)
		return
	}

	_, err = executeScriptSync(d, string(c), args)
	if err != nil {
		log.Println("error on execute script", err)
		return
	}

	return
}


func (d *Driver) FindElementByXpathJs(args ...interface{}) interface{} {
	f := file.FindFile(config.TestSetting.JsFilesPath, "findElementByXpath.js")

	c, err := readFile(f)
	if err != nil {
		log.Println("error on file read in click.js", err)
		return nil
	}

	rtn, err := executeScriptSync(d, string(c), args)
	if err != nil {
		log.Println("error on execute script", err)
		return nil
	}

	return rtn
}

func (d *Driver) ExecuteScript(fName string, args ...interface{}) interface{} {
	f := file.FindFile(config.TestSetting.JsFilesPath, fName)

	c, err := os.ReadFile(f)
	if err != nil {
		log.Println("error on read file", err)
		return nil
	}

	rtn, err := executeScriptSync(d, string(c), args)
	if err != nil {
		log.Println("error on execute script", err)
		return nil
	}

	return rtn
}

func (d *Driver) Script(script string, args ...interface{}) interface{} {
	res, err := executeScriptSync(d, script, args...)
	if err != nil {
		log.Println("error execute script:", err)
		return nil
	}

	return res
}

// executeScriptSync
func executeScriptSync(d *Driver, script string, args ...interface{}) (interface{}, error) {
	if args == nil {
		args = make([]interface{}, 0)
	}

	rr := new(struct{ Value interface{} })
	d.execute(defaultStrategy{Command{
		Path:   PathDriverScriptSync,
		Method: http.MethodPost,
		Data: marshalData(map[string]interface{}{
			"script": script,
			"args":   args,
		}),
		ResponseData: rr,
	}})

	return rr.Value, nil
}
