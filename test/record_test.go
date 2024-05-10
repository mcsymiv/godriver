package test

import (
	"testing"

	"github.com/mcsymiv/godriver/record"
)

func TestRecord(t *testing.T) {
	record.CreateTest()
}

func TestSteps(t *testing.T) {
	record.CreateSteps()
}
