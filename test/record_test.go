package test

import (
	"testing"

	"github.com/mcsymiv/godriver/record"
)

func TestSteps(t *testing.T) {
	var tName string = "../test/steps_test.go"
	var rName string = "record_3.json"

	record.CreateSteps(tName, rName)
}
