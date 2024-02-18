package driver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// IsDisplayed
func (e *Element) IsDisplayed() *Element {
	dis, err := isDisplayed(e)
	if err != nil {
		log.Println("error on displayed")
		return nil
	}

	if !dis {
		log.Println("element not visible")
		return nil
	}

	return e
}

type displayStrategy struct {
	findReq *http.Request
	*Driver
	timeout, delay time.Duration
}

// newDisplayCommand
// referrences find elememnt cmd
func newDisplayCommand(e *Element) *Command {
	fCommand := newFindCommand(e.Selector, e.Driver)
	fReq, _ := newCommandRequest(e.Client, fCommand)

	return &Command{
		Path:           "/element/%s/displayed",
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodGet,

		Strategies: []CommandExecutor{
			&displayStrategy{
				Driver:  e.Driver,
				findReq: fReq,
				timeout: 8,
				delay:   800,
			},
		},
	}
}

func (dis displayStrategy) Execute(req *http.Request) (*http.Response, error) {
	var displayRes = new(struct{ Value bool })
	var buffRes *buffResponse
	var res *http.Response
	var err error

	// perform isDisplayed check
	res, err = dis.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// convert response to NopCloser
	buffRes = newBuffResponse(res)
	unmarshalData(buffRes.Response, displayRes)

	// start waiter check
	if !displayRes.Value {
		log.Println("element is not visible")
		start := time.Now()
		end := start.Add(dis.timeout * time.Second)

		for {
			log.Println("element still not visible")
			time.Sleep(dis.delay * time.Millisecond)
			res, err = dis.Client.HTTPClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("error")
			}

			buffRes = newBuffResponse(res)
			unmarshalData(buffRes.Response, displayRes)

			if displayRes.Value {
				return buffRes.Response, nil
			}

			if time.Now().After(end) {
				dis.Screenshot()
				return res, fmt.Errorf("error on element display timeout: %v", err)
			}
		}
	}

	return buffRes.Response, err
}

func isDisplayed(e *Element) (bool, error) {
	op := newDisplayCommand(e)

	res := e.Client.ExecuteCommand(op)
	d := new(struct{ Value bool })
	unmarshalData(res[0].Response, d)

	return d.Value, nil
}
