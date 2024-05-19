package driver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/mcsymiv/godriver/capabilities"
	"github.com/mcsymiv/godriver/config"
)

const (
	PathDriverUrl           = "/url"
	PathDriverFrame         = "/frame"
	PathDriverRefresh       = "/refresh"
	PathDriverWindow        = "/windown"
	PathDriverWindowNew     = "/windown/new"
	PathDriverWindowHandles = "/windown/handles"

	PathElementFind         = "/element"
	PathElementsFind        = "/elements"
	PathElementActive       = "/element/active"
	PathElementDisplayed    = "/element/%s/displayed"
	PathElementClear        = "/element/%s/clear"
	PathElementValue        = "/element/%s/value"
	PathElementFromElement  = "/element/%s/element"
	PathElementsFromElement = "/element/%s/elements"
	PathElementAttribute    = "/element/%s/attribute/%s"
	PathElementClick        = "/element/%s/click"
	PathElementText         = "/element/%s/text"
)

type Driver struct {
	Client       *Client
	Session      *Session
	ServiceCmd   *exec.Cmd
	Capabilities *capabilities.Capabilities
}

// NewDriver
// Webdriver setup
// 1. starts webdriver service based on browser name capabilities
// 2. wait for service process to start. requests /status with 2 second timeuout
// 3. creates new session to use
// 4. initializes new client
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
		fmt.Sprintf("%s://%s:%s/session/%s",
			caps.DriverSetupCapabilities.Protocol,
			caps.DriverSetupCapabilities.Host,
			caps.DriverSetupCapabilities.Port,
			s.Id,
		),
	)

	config.TestSetting = config.DefaultSetting()

	return &Driver{
		Client:       c,
		ServiceCmd:   cmd,
		Capabilities: &caps,
	}
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
