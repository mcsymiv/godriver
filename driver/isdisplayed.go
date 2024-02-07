package driver

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

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

func isDisplayed(e *Element) (bool, error) {
	op := &Command{
		Path:   fmt.Sprintf("/element/%s/displayed", e.Id),
		Method: http.MethodGet,
	}

	st := &displayStrategy{
		Client:  *e.Client.HTTPClient,
		Timeout: 15,
		Delay:   1000,
	}

	res, err := e.Client.ExecuteCommandStrategy(op, st)
	if err != nil {
		return false, err
	}

	d := new(struct{ Value bool })
	unmarshalData(res, d)

	return d.Value, nil
}

type displayStrategy struct {
	http.Client
	Timeout, Delay time.Duration
}

type BuffResponse struct {
	*http.Response
	bBuffer *bytes.Buffer
}

func newBuffResponse(response *http.Response) *BuffResponse {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("error on buf write")
		return nil
	}

	buffRes := &BuffResponse{
		Response: response,
		bBuffer:  bytes.NewBuffer(body),
	}

	rr := io.LimitReader(ReusableReader(buffRes.bBuffer), 2048*2)
	resBody := io.NopCloser(rr)
	buffRes.Response.Body = resBody

	return buffRes
}

func (dis *displayStrategy) Execute(req *http.Request) (*http.Response, error) {
	log.Printf("displaye strategy request: %s", req.URL.Path)
	var start time.Time = time.Now()
	var end time.Time = start.Add(dis.Timeout * time.Second)

	var displayRes = new(struct{ Value bool })
	var dBuffRes *BuffResponse
	var res *http.Response
	var err error

	res, err = dis.Client.Do(req)

	for {
		if err != nil {
			log.Printf("element not fount: %v", res.StatusCode)
			log.Printf("error in display: %v", err)
		}

		dBuffRes = newBuffResponse(res)
		unmarshalData(dBuffRes.Response, displayRes)

		if displayRes.Value {
			log.Printf("element is visible in execute: %v", displayRes)
			return res, err
		}

		time.Sleep(dis.Delay * time.Millisecond)
		res, err := dis.Client.Do(req)

		if time.Now().After(end) {
			log.Printf("timeout")
			return res, err
		}
	}
}
