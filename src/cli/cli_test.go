package cli

import (
	"fmt"
	"testing"
)

type Read struct {
	returnMessage string
}

var Iter = 0
var Done = false

func (read *Read) ReadString(s string) string {
	return read.returnMessage
}

func TestExit(t *testing.T) {
	if !shouldExit(func() string {
		return "y"
	}) {
		t.Fail()
	}
	if !shouldExit(func() string {
		return "yes"
	}) {
		t.Fail()
	}
	if !shouldExit(func() string {
		return "Y"
	}) {
		t.Fail()
	}
	if !shouldExit(func() string {
		return "Yes"
	}) {
		t.Fail()
	}
	if shouldExit(func() string {
		return "y "
	}) {
		t.Fail()
	}
	if shouldExit(func() string {
		return "n"
	}) {
		t.Fail()
	}
	if shouldExit(func() string {
		return "N"
	}) {
		t.Fail()
	}
}

func TestTextStringCompare(t *testing.T) {
	if !stringsEqual("TestText", "TestText") {
		t.Fail()
	}
	if stringsEqual("TestText", "TestText2") {
		t.Fail()
	}
	if stringsEqual("TestText", "testText") {
		t.Fail()
	}
}

func TestPrintHelp(t *testing.T) {
	printHelp()
}

func TestDoExit(t *testing.T) {
	Iter = 0
	Done = false
	do(func() string {
		if Iter == 0 {
			Iter++
			fmt.Println("Returning exit")
			return "exit"
		} else if Iter == 1 {
			fmt.Println("Returning y")
			return "y"
		}
		return ""
	}, func() {
		Done = true
	})
	if !Done {
		t.Fail()
	}
}