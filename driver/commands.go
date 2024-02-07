package driver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type WrapCommand func(CommandHandler) CommandHandler

type CommandHandler func(*http.Request) (*http.Response, error)

// Route represents a specific API route and its handler function.
type Command struct {
	Path, Method string
	Data         []byte
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

func marshalData(body interface{}) []byte {
	b, err := json.Marshal(body)
	if err != nil {
		log.Println("error on marshal: ", err)
		return nil
	}

	return b
}

func unmarshalData(res *http.Response, any interface{}) []byte {
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error on reading response:", err)
		return nil
	}

	if err := json.Unmarshal(b, &any); err != nil {
		log.Println("error on unmarshal:", err)
		return nil
	}

	return b
}
