package driver

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/mcsymiv/godriver/capabilities"
)

var GeckoDriverPath string = "geckodriver"

// ubuntu version
// var ChromeDriverPath string = "chromerdriver"

// mac version
var ChromeDriverPath string = "/Users/mcs/tools/chromedriver_122"

type DriverStatus struct {
	Message string `json:"message"`
	Ready   bool   `json:"ready"`
}

// TODO: add service wrapper for log file and exec.Cmd
var OutFileLogs *os.File

// newService
// starts local webdriver service, geckodriver or chromedriver
// based on passed capabilities i.e. cap.Capabilities.AlwaysMatch.BrowserName
// defaults to firefox
func newService(caps *capabilities.Capabilities) (*exec.Cmd, error) {
	// returns command arguments for specified driver to start from shell
	var cmdArgs []string = driverCommand(caps)

	// previously used line to start driver
	// cmd := exec.Command("zsh", "-c", GeckoDriverrequest, "--port", "4444", ">", "logs/gecko.session.logs", "2>&1", "&")
	// open the out file for writing
	OutFileLogs, err := os.Create("../artifacts/logs/logs.txt")
	if err != nil {
		log.Println("failed to start driver service:", err)
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = OutFileLogs
	cmd.Stderr = OutFileLogs

	err = cmd.Start()
	if err != nil {
		log.Println("failed to start driver service:", err)
		return nil, err
	}

	if cmd.Process.Pid == 0 {
		return nil, fmt.Errorf("service did not start")
	}

	return cmd, nil
}

// driverCommand
// Check for specified driver/browser name to pass to cmd to start the driver server
func driverCommand(cap *capabilities.Capabilities) []string {
	// when calling /bin/zsh -c command
	// command arguments will be ignored
	var cmdArgs []string = []string{
		// "-c",
	}

	if cap.Capabilities.AlwaysMatch.BrowserName == "firefox" {
		cmdArgs = append(cmdArgs, GeckoDriverPath, "--port", cap.Port, "--log", "trace")
	} else {
		cmdArgs = append(cmdArgs, ChromeDriverPath, fmt.Sprintf("--port=%s", cap.Port), "--verbose", "--whitelisted-ips", "--log-path=chromedriver.log", "--enable-chrome-browser-cloud-management")
		// cmdArgs = append(cmdArgs, ChromeDriverPath, fmt.Sprintf("--port=%s", cap.Port))
	}

	// redirect output argumetns ignored when used in exec.Command
	// cmdArgs = append(cmdArgs, ">", "logs/session.log", "2>&1", "&")
	return cmdArgs
}

func waitForDriverService(cmd *exec.Cmd, caps *capabilities.Capabilities) error {
	// Tries to get driver status for 2 seconds
	// Once driver isReady, returns command for deferred kill
	start := time.Now()
	end := start.Add(4 * time.Second)
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
