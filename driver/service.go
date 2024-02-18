package driver

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/mcsymiv/godriver/capabilities"
)

var GeckoDriverPath string = "/Users/mcs/Documents/tools/geckodriver"
var ChromeDriverPath string = "/Users/mcs/Documents/tools/chromedriver"

type DriverStatus struct {
	Message string `json:"message"`
	Ready   bool   `json:"ready"`
}

// newService
// starts local webdriver service, geckodriver or chromedriver
// based on passed capabilities i.e. cap.Capabilities.AlwaysMatch.BrowserName
// defaults to firefox
func newService(caps *capabilities.Capabilities) (*exec.Cmd, error) {
	// returns command arguments for specified driver to start from shell
	var cmdArgs []string = driverCommand(caps)

	// previously used line to start driver
	// cmd := exec.Command("zsh", "-c", GeckoDriverrequest, "--port", "4444", ">", "logs/gecko.session.logs", "2>&1", "&")
	cmd := exec.Command("/bin/zsh", cmdArgs...)
	err := cmd.Start()
	if err != nil {
		log.Println("failed to start driver service:", err)
		return nil, err
	}

	return cmd, nil
}

// driverCommand
// Check for specified driver/browser name to pass to cmd to start the driver server
func driverCommand(cap *capabilities.Capabilities) []string {
	var cmdArgs []string = []string{
		"-c",
	}

	if cap.Capabilities.AlwaysMatch.BrowserName == "firefox" {
		cmdArgs = append(cmdArgs, GeckoDriverPath, "--port", cap.Port)
	} else {
		cmdArgs = append(cmdArgs, ChromeDriverPath, fmt.Sprintf("--port=%s", cap.Port))
	}

	cmdArgs = append(cmdArgs, ">", "logs/session.log", "2>&1", "&")
	return cmdArgs
}

func waitForDriverService(cmd *exec.Cmd, caps *capabilities.Capabilities) error {
	// Tries to get driver status for 2 seconds
	// Once driver isReady, returns command for deferred kill
	start := time.Now()
	end := start.Add(2 * time.Second)
	for stat, err := getDriverStatus(caps); err != nil || !stat.Ready; stat, err = getDriverStatus(caps) {
		time.Sleep(200 * time.Millisecond)
		log.Println("Error getting driver status:", err)

		if time.Now().After(end) {
			log.Println("Killing cmd:", cmd)
			cmd.Process.Kill()
			return errors.New("driver start timed out after 2 seconds")
		}
	}

	return nil
}
