package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mcsymiv/godriver/capabilities"
)

type JsonSession struct {
	Id string `json:"sessionId"`
}

type Session struct {
	Id, Route string
}

func newSession(caps *capabilities.Capabilities) (*Session, error) {
	data, err := json.Marshal(caps)
	if err != nil {
		log.Printf("new driver marshall error: %+v", err)
		return nil, err
	}
	url := fmt.Sprintf("http://%s:%s/session", caps.Host, caps.Port)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "json/application")

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	reply := new(struct{ Value JsonSession })
	if err := json.Unmarshal(body, reply); err != nil {
		log.Println("Status unmarshal error", err)
		return nil, err
	}

	log.Println(reply.Value)

	return &Session{
		Id:    reply.Value.Id,
		Route: "/session",
	}, nil
}

func (d Driver) Quit() {
	q := &Command{
		Path:   "",
		Method: http.MethodDelete,
	}

	d.Client.ExecuteCommandStrategy(q)
}
