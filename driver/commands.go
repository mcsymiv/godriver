package driver

import (
	"bytes"
	"fmt"
	"net/http"
)

// Command
// represents a request (command) to webdriver
// with added Strategies to execute in CommandExecutor
type Command struct {
	Path           string
	Method         string
	PathFormatArgs []any

	Data         []byte
	ResponseData interface{}

	Strategy CommandExecutor
}

// CommandExecutor
// strategy to remove duplicates in execute Command/Request
type CommandExecutor interface {

	// TODO: Exec wrapper around req/res
	// Exec(r *buffRequest) (*buffResponse, error)

	exec(*Command, interface{})
}

// Context
type CommandStrategy struct {
	CommandExecutor
}

type execContext struct {
	cmd CommandExecutor
}

// newCommandRequest
// updated version without Session
func newCommandRequest(c *Client, cmd *Command) (*http.Request, error) {
	var cPath string = cmd.Path
	if len(cmd.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(cmd.Path, cmd.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, cPath)

	// TODO: investigate on chrome API, read upon bytes.NewReader vs bytes.NewBuffer diff
	// review if reusable request is needed(?)
	// chromedriver does not accept reusable Reader, i.e. bytes.NewReader(cmd.Data)
	// but, bytes.NewBuffer() without NopCloser wrappper works(!)
	// note: geckodriver is fine with reusable request reader

	// rUse := ReusableReader(bytes.NewBuffer(cmd.Data))
	// rr := io.LimitReader(rUse, c.RequestReaderLimit)
	// reqBody := io.NopCloser(rr)
	// req, err := http.NewRequest(cmd.Method, url, reqBody)

	req, err := http.NewRequest(cmd.Method, url, bytes.NewBuffer(cmd.Data))
	if err != nil {
		return nil, fmt.Errorf("error on new request: %v", err)
	}

	return req, nil
}
