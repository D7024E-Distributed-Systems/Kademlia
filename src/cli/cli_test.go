package cli

import (
	"fmt"
	"testing"

	"github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
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
	contact := kademlia.NewContact(kademlia.NewRandomKademliaID(), "localhost")
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

	}, kademlia.NewKademliaStruct(kademlia.NewNetwork(&contact)))
	if !Done {
		t.Fail()
	}
}

func TestFindContact(t *testing.T) {
	findContact(func() string {
		return "A000000000000000000000000000000000000000"
	}, func(kad *kademlia.KademliaID) *kademlia.Contact {
		res := kademlia.NewContact(kad, "localhost")
		return &res
	})
}

func TestFindContactNil(t *testing.T) {
	findContact(func() string {
		return "test"
	}, func(kad *kademlia.KademliaID) *kademlia.Contact {
		return nil
	})
}
