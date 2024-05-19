package test

import (
	"testing"

	"github.com/mcsymiv/godriver/record"
)

func TestSteps(t *testing.T) {
	var fName string = "../test/brn_test.go"
	var rName string = "brn.json"
	var tName string = "Brn"

	record.CreateSteps(fName, rName, tName)
}
