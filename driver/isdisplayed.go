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

	// get response buffer
	// reads response body
	buffRes = newBuffResponse(res)
	unmarshalResponses([]*buffResponse{buffRes}, displayRes)

	// start waiter check
	if !displayRes.Value {
		start := time.Now()
		end := start.Add(dis.timeout * time.Second)

		for {
			log.Println("element is not visible")
			time.Sleep(dis.delay * time.Millisecond)
			res, err = dis.Client.HTTPClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("error")
			}

			buffRes = newBuffResponse(res)
			unmarshalResponses([]*buffResponse{buffRes}, displayRes)

			if displayRes.Value {
				// set NopCloser response with body
				buffRes.Response.Body = buffRes.bRead()
				return buffRes.Response, nil
			}

			if time.Now().After(end) {
				dis.Screenshot()
				return buffRes.Response, fmt.Errorf("error on element display timeout: %v", err)
			}
		}
	}

	// set NopCloser response with body
	buffRes.Response.Body = buffRes.bRead()
	return buffRes.Response, err
}

func isDisplayed(e *Element) (bool, error) {
	op := newDisplayCommand(e)

	d := new(struct{ Value bool })
	e.Client.ExecuteCmd(op, d)

	return d.Value, nil
}
