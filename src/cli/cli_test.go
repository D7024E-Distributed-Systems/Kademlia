package cli

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
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

func TestFindContactWithDo(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&contact)
	quit := make(chan bool)
	i := 0
	go do(func() string {
		if i == 0 {
			i++
			return "find contact"
		} else {
			quit <- true
			return ""
		}
	}, func() {
		quit <- true
	}, NewKademliaStruct(network))
}

func TestPutWithDo(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&contact)
	quit := make(chan bool)
	i := 0
	go do(func() string {
		if i == 0 {
			i++
			return "put"
		} else {
			quit <- true
			return ""
		}
	}, func() {
		quit <- true
	}, NewKademliaStruct(network))
}

func TestHelpWithDo(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&contact)
	quit := make(chan bool)
	i := 0
	go do(func() string {
		if i == 0 {
			i++
			return "help"
		} else {
			quit <- true
			return ""
		}
	}, func() {
		quit <- true
	}, NewKademliaStruct(network))
}

func TestForgetWithDo(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&contact)
	quit := make(chan bool)
	i := 0
	go do(func() string {
		if i == 0 {
			i++
			return "forget"
		} else {
			quit <- true
			return ""
		}
	}, func() {
		quit <- true
	}, NewKademliaStruct(network))
}

func TestGetWithDo(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&contact)
	quit := make(chan bool)
	i := 0
	go do(func() string {
		if i == 0 {
			i++
			return "get"
		} else {
			quit <- true
			return ""
		}
	}, func() {
		quit <- true
	}, NewKademliaStruct(network))
}
func TestNothingWithDo(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&contact)
	quit := make(chan bool)
	i := 0
	go do(func() string {
		if i == 0 {
			i++
			return "abc"
		} else {
			quit <- true
			return ""
		}
	}, func() {
		quit <- true
	}, NewKademliaStruct(network))
}

func TestGetTableDo(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&contact)
	quit := make(chan bool)
	i := 0
	go do(func() string {
		if i == 0 {
			i++
			return "table"
		} else {
			quit <- true
			return ""
		}
	}, func() {
		quit <- true
	}, NewKademliaStruct(network))
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
	}, func(kad *kademlia.KademliaID) ContactCandidates {
		candidates := ContactCandidates{}
		res := kademlia.NewContact(kad, "localhost")
		res2 := kademlia.NewContact(kad, "localhost")
		contacts := []kademlia.Contact{res, res2}
		candidates.Append(contacts)
		return candidates
	})
}

func TestFindContact2(t *testing.T) {
	findContact(func() string {
		return "A0000000000000000000000000000000000000000"
	}, func(kad *kademlia.KademliaID) ContactCandidates {
		candidates := ContactCandidates{}
		res := kademlia.NewContact(kad, "localhost")
		res2 := kademlia.NewContact(kad, "localhost")
		contacts := []kademlia.Contact{res, res2}
		candidates.Append(contacts)
		return candidates
	})
}

func TestForgetSuccess(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost:3000")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)
	kademlia.KnownHolders[contact] = *ToKademliaID("A000000000000000000000000000000000000000")
	forgetHelp(kademlia, func() string {
		return "A000000000000000000000000000000000000000"
	}, func(value string) bool {
		if value == "A000000000000000000000000000000000000000" {
			return true
		}
		return false
	})
}

func TestForgetFail(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost:3000")
	kademlia := NewKademliaStruct(NewNetwork(&contact))
	kademlia.KnownHolders[contact] = *ToKademliaID("B000000000000000000000000000000000000000")
	forgetHelp(kademlia, func() string {
		return "A000000000000000000000000000000000000000"
	}, func(value string) bool {
		if value != "A000000000000000000000000000000000000000" {
			return true
		}
		return false
	})
}

func TestGetSuccess(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost:3000")
	kademlia := NewKademliaStruct(NewNetwork(&contact))
	kademlia.Store([]byte("File"), 15*time.Minute)
	id := NewKademliaID("File")
	stringId := id.String()
	val := kademlia.LookupData(*id)
	if string(val) != "File" {
		fmt.Println(val)
		t.Fail()
	}
	getValue(
		func() string { return stringId },
		kademlia,
	)
}

func TestGetFailure(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost:3000")
	kademlia := NewKademliaStruct(NewNetwork(&contact))
	id := NewKademliaID("File")
	stringId := id.String()
	getValue(
		func() string { return stringId },
		kademlia,
	)
}

func TestGetFailure2(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost:3000")
	kademlia := NewKademliaStruct(NewNetwork(&contact))
	getValue(
		func() string { return "A0000000000000000000000000000000000000000" },
		kademlia,
	)
}

func TestStoreValue(t *testing.T) {
	failed := 0
	storeValue(func() string {
		return "10m"
	}, func(data []byte, t time.Duration) ([]*KademliaID, string) {
		if bytes.Compare(data, []byte("10m")) != 0 {
			failed += 1
		}
		if t != time.Minute*10 {
			failed += 1
		}
		return []*KademliaID{NewRandomKademliaID()}, ""
	})
	if failed > 0 {

		t.Fail()
	}
}

func TestStoreValueFailure(t *testing.T) {
	failed := 0
	go storeValue(func() string {
		return "10x"
	}, func(data []byte, t time.Duration) ([]*KademliaID, string) {
		if bytes.Compare(data, []byte("10x")) != 0 {
			failed += 1
		}
		if t != time.Minute*10 {
			failed += 1
		}
		return []*KademliaID{NewRandomKademliaID()}, ""
	})
	time.Sleep(1 * time.Millisecond)
	return

}

func TestInit(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&contact)
	go Init(func() {
		os.Exit(0)
	}, NewKademliaStruct(network))
	time.Sleep(1 * time.Millisecond)
	return
}
