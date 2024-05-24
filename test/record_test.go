package test

import (
	"testing"

	"github.com/mcsymiv/godriver/config"
	"github.com/mcsymiv/godriver/record"
)

func TestSteps(t *testing.T) {
	config.TestSetting = config.DefaultSetting()

	var fName string = "../test/new_acc_test.go"
	var rName string = "new_account.json"
	var tName string = "NewAcc"

	record.CreateSteps(fName, rName, tName)
}
