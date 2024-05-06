package test

import (
	"testing"

	"github.com/mcsymiv/godriver/config"
)

func TestRecord(t *testing.T) {
	config.CreateTest()
}

func TestSteps(t *testing.T) {
	config.CreateSteps()
}
