package driver

import (
	"log"
	"net/http"
	"os"

	"github.com/mcsymiv/godriver/config"
)

func (d Driver) ExecuteScript(fName string, args ...string) interface{} {
	f := config.FindFile("../js", fName)

	// only one arg will be applied
	// slice serves as "optional" argument
	if len(args) == 1 {
		for _, arg := range args {
			act := &config.ReplaceWord{
				ReplaceLine: &config.ReplaceLine{
					Old: "<placeholder>",
					New: arg,
				},
			}
			config.Exec(act, f)
		}
	}

	c, err := os.ReadFile(f)
	if err != nil {
		log.Println("error on read file", err)
		return nil
	}

	rtn, err := executeScriptSync(d, string(c))
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

func executeScriptSync(d Driver, script string, args ...interface{}) (interface{}, error) {
	if args == nil {
		args = make([]interface{}, 0)
	}

	op := &Command{
		Path:   "/execute/sync",
		Method: http.MethodPost,
		Data: marshalData(map[string]interface{}{
			"script": script,
			"args":   args,
		}),
	}

	bRes := d.Client.ExecuteCommand(op)
	rr := new(struct{ Value interface{} })
	unmarshalResponses(bRes, rr)

	return rr.Value, nil
}
