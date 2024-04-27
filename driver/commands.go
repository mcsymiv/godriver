package driver

import (
	"bytes"
	"encoding/json"
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
	Exec(r *buffRequest) (*buffResponse, error)
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

// ExecuteCmd
//  1. general purpose client receiver
//     executes prepared command and strategies (if defined)
//     when no strategy difened, executes client request
//  2. unmarshals passed data struct
//
// TODO: refactor to internal use, i.e. executeCmd
func (c *Client) ExecuteCmd(cmd *Command, d ...any) ([]*buffResponse, error) {
	req, err := newCommandRequest(c, cmd)
	if err != nil {
		return nil, fmt.Errorf("error on new command request: %v", err)
	}

	st := newExecutorContext(c, cmd)
	for i, s := range st.cmds {

		// executes request inside defined CommandExecutor strategy
		// if none provided, performs http.Request with Client's DefaultExecuteStrategy
		res, err := NewStrategy(s).Execute(req)
		if err != nil {
			return nil, fmt.Errorf("error on new strategy exec: %+v", err)
		}

		st.bufs[i], err = newBuffResponse(res)
		if err != nil {
			return nil, fmt.Errorf("error on new buffered response: %v", err)
		}
	}

	if len(st.bufs) > 0 && len(d) > 0 {
		for i, res := range st.bufs {
			err := json.Unmarshal(res.buff, d[i])

			if err != nil {
				return nil, fmt.Errorf("error on unmarshal %d response: %v", i, err)
			}
		}
	}

	return st.bufs, nil
}
