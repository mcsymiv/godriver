package driver

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// type WrapCommand func(CommandHandler) CommandHandler

// type CommandHandler func(*http.Request) (*http.Response, error)

// Route represents a specific API route and its handler function.
type Command struct {
	Path           string
	Method         string
	PathFormatArgs []any

	Data         []byte
	ResponseData interface{}

	Strategies []CommandExecutor
}

// strategy to remove duplicates in execute Command/Request
type CommandExecutor interface {
	Execute(req *http.Request) (*http.Response, error)
	Exec(r *buffRequest) (*buffResponse, error) // TODO: Exec wrapper around req/res
}

// Context
type CommandStrategy struct {
	CommandExecutor
}

type buffResponse struct {
	*http.Response
	buff  []byte
	bRead func() io.ReadCloser // callback when response with body required
}

type buffRequest struct {
	*http.Request
	bRead func() io.ReadCloser
}

type executorContext struct {
	cmds []CommandExecutor
	bufs []*buffResponse
}

// newExecutorContext
// creates new CommandExecutor out of defaul client
// or from passed in command strategies
// allocates space for buffered command response
func newExecutorContext(c *Client, cmd *Command) *executorContext {
	var executorCtx *executorContext = &executorContext{
		cmds: []CommandExecutor{c},
		bufs: make([]*buffResponse, 1),
	}

	if len(cmd.Strategies) > 0 {
		executorCtx.cmds = cmd.Strategies
		executorCtx.bufs = make([]*buffResponse, len(cmd.Strategies))
	}

	return executorCtx
}

// newBuffResponse
// reusable response for multiple reads
func newBuffResponse(response *http.Response) (*buffResponse, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error on read all body response: %v", err)
	}

	buffRes := &buffResponse{
		buff:     body,
		Response: response,
		bRead: func() io.ReadCloser {
			rr := io.LimitReader(ReusableReader(bytes.NewBuffer(body)), 2048*2)
			return io.NopCloser(rr)
		},
	}

	return buffRes, nil
}

// newCommandRequest
// updated version without Session
func newCommandRequest(c *Client, cmd *Command) (*http.Request, error) {
	var cPath string = cmd.Path
	if len(cmd.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(cmd.Path, cmd.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, cPath)

	rr := io.LimitReader(ReusableReader(bytes.NewReader(cmd.Data)), c.RequestReaderLimit)
	reqBody := io.NopCloser(rr)
	req, err := http.NewRequest(cmd.Method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error on new request: %v", err)
	}

	return req, nil
}
