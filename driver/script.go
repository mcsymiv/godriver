package driver

import (
	"log"
	"net/http"
	"os"

	"github.com/mcsymiv/godriver/config"
	"github.com/mcsymiv/godriver/file"
)

func (d Driver) ExecuteScript(fName string, args ...interface{}) interface{} {
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

func (d Driver) Script(script string, args ...interface{}) interface{} {
	res, err := executeScriptSync(d, script, args...)
	if err != nil {
		log.Println("error execute script:", err)
		return nil
	}

	return res
}

// executeScriptSync
func executeScriptSync(d Driver, script string, args ...interface{}) (interface{}, error) {
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
