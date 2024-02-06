package driver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/mcsymiv/godriver/capabilities"
)

type Driver struct {
	Client       *Client
	Session      *Session
	ServiceCmd   *exec.Cmd
	Commands     map[string]*Command
	Capabilities *capabilities.Capabilities
}

func NewDriver(capsFn ...capabilities.CapabilitiesFunc) *Driver {
	caps := capabilities.DefaultCapabilities()
	for _, capFn := range capsFn {
		capFn(&caps)
	}

	// newService
	// calls driver binary local path
	// starts driver service with provided capabilities
	cmd, err := newService(&caps)
	if err != nil {
		log.Fatal("unable to start driver service", err)
	}

	// waitForDriverService
	// there is a delay in geckodriver responce to /status request
	// waits for driver to be ready to accept incoming requests
	err = waitForDriverService(cmd, &caps)
	if err != nil {
		log.Fatal("driver start timed out", err)
	}

	// newSession
	// return session id to use
	s, err := newSession(&caps)
	if err != nil || s == nil {
		log.Fatal("unable to start session", err)
	}

	// newClient
	// utilizes RoundTripper interface to wrap requests to the driver
	c := newClient(
		fmt.Sprintf("%s://%s:%s", caps.DriverSetupCapabilities.Protocol, caps.DriverSetupCapabilities.Host, caps.DriverSetupCapabilities.Port),
		s,
	)

	dCommands := registerCommands()

	return &Driver{
		Client:       c,
		Session:      s,
		ServiceCmd:   cmd,
		Capabilities: &caps,
		Commands:     dCommands,
	}
}

func registerCommands() map[string]*Command {
	var cmds = make(map[string]*Command)
	cmds["open"] = &Command{Path: "/url", Method: http.MethodPost}
	cmds["refresh"] = &Command{Path: "/refresh", Method: http.MethodPost}
	cmds["find"] = &Command{Path: "/element", Method: http.MethodPost}
	cmds["finds"] = &Command{Path: "/elements", Method: http.MethodPost}
	cmds["frame"] = &Command{Path: "/frame", Method: http.MethodPost}

	return cmds
}

// Service
// Returns ref to started driver service process
func (d Driver) Service() *exec.Cmd {
	return d.ServiceCmd
}

func getDriverStatus(caps *capabilities.Capabilities) (*DriverStatus, error) {
	url := fmt.Sprintf("%s://%s:%s/status", caps.DriverSetupCapabilities.Protocol, caps.DriverSetupCapabilities.Host, caps.DriverSetupCapabilities.Port)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	reply := new(struct{ Value DriverStatus })
	if err := json.Unmarshal(body, reply); err != nil {
		log.Println("Status unmarshal error", err)
		return &reply.Value, err
	}

	return &reply.Value, nil
}
