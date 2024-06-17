package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mcsymiv/godriver/capabilities"
)

type JsonSession struct {
	Id string `json:"sessionId"`
}

type Session struct {
	Id, Route string
}

// newSession
// hardcoded request to start new webdriver session
func newSession(caps *capabilities.Capabilities) (*Session, error) {
	data, err := json.Marshal(caps)
	if err != nil {
		return nil, fmt.Errorf("error on new driver marshal caps: %v", err)
	}
	url := fmt.Sprintf("http://%s:%s/session", caps.Host, caps.Port)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error on session request: %v", err)
	}

	req.Header.Add("Accept", "json/application")

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on create session request: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error on read responce body")
	}

	reply := new(struct{ Value JsonSession })
	if err := json.Unmarshal(body, reply); err != nil {
		return nil, fmt.Errorf("erron on status unmarshal: %v", err)
	}

	return &Session{
		Id:    reply.Value.Id,
		Route: "/session",
	}, nil
}

// Quit
// closes active webdriver session
func (d *Driver) Quit() {
	d.execute(defaultStrategy{Command{
		Path:   "",
		Method: http.MethodDelete,
	}})
}
