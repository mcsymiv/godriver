package driver

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type WrapCommand func(CommandHandler) CommandHandler

type CommandHandler func(*http.Request) (*http.Response, error)

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
}

// Context
type CommandStrategy struct {
	CommandExecutor
}

func NewStrategy(cmd CommandExecutor) *CommandStrategy {
	return &CommandStrategy{
		CommandExecutor: cmd,
	}
}

func (c *CommandStrategy) Set(cmd CommandExecutor) {
	c.CommandExecutor = cmd
}

func (c CommandStrategy) Exec(req *http.Request) (*http.Response, error) {
	return c.CommandExecutor.Execute(req)
}

type buffResponse struct {
	*http.Response
	buff  []byte
	bRead func() io.ReadCloser // callback when response with body required
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
	var st []CommandExecutor
	var buffRes []*buffResponse

	if len(cmd.Strategies) == 0 {
		st = []CommandExecutor{c}
		buffRes = make([]*buffResponse, 1)
	} else {
		st = cmd.Strategies
		buffRes = make([]*buffResponse, len(cmd.Strategies))
	}

	return &executorContext{
		cmds: st,
		bufs: buffRes,
	}
}

// newBuffResponse
// reusable response for multiple reads
// TODO: think about passing in interface{} of returned struct
func newBuffResponse(response *http.Response) *buffResponse {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("error on buf write")
		return nil
	}

	buffRes := &buffResponse{
		buff:     body,
		Response: response,
		bRead: func() io.ReadCloser {
			rr := io.LimitReader(ReusableReader(bytes.NewBuffer(body)), 2048*2)
			return io.NopCloser(rr)
		},
	}

	return buffRes
}

// newCommandRequest
func newCommandRequest(c *Client, cmd *Command) (*http.Request, error) {
	var cPath string = cmd.Path
	if len(cmd.PathFormatArgs) != 0 {
		cPath = fmt.Sprintf(cmd.Path, cmd.PathFormatArgs...)
	}

	url := fmt.Sprintf("%s%s/%s%s", c.BaseURL, c.Session.Route, c.Session.Id, cPath)

	rr := io.LimitReader(ReusableReader(bytes.NewReader(cmd.Data)), c.RequestReaderLimit)
	reqBody := io.NopCloser(rr)
	req, err := http.NewRequest(cmd.Method, url, reqBody)
	if err != nil {
		return nil, err
	}

	return req, nil
}
