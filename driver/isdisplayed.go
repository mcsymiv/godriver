package driver

import (
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
	*http.Client
	timeout, delay time.Duration
}

// newDisplayCommand
// referrences find elememnt cmd
func newDisplayCommand(e *Element) *Command {
	fCommand := newFindCommand(e.By, e.Driver)
	fReq, _ := newCommandRequest(e.Client, fCommand)

	return &Command{
		Path:           "/element/%s/displayed",
		PathFormatArgs: []any{e.Id},
		Method:         http.MethodGet,

		Strategies: []CommandExecutor{
			&displayStrategy{
				Client:  e.Client.HTTPClient,
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
	res, err = dis.Client.Do(req)
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
			time.Sleep(dis.delay * time.Millisecond)
			res, err = dis.Client.Do(req)
			if err != nil {
				return nil, err
			}

			buffRes = newBuffResponse(res)
			unmarshalData(buffRes.Response, displayRes)

			if displayRes.Value {
				return buffRes.Response, nil
			}

			if time.Now().After(end) {
				log.Printf("timeout")
				return res, err
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
