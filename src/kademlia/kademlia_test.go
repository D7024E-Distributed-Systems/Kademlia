package kademlia

import (
	"fmt"
	"testing"
)

func TestInsertData(t *testing.T) {
	kd := NewKademliaStruct()
	if len(kd.m) != 0 {
		fmt.Println(len(kd.m))
		t.Fail()
	}

	kd.Store([]byte("AA"))

	if len(kd.m) != 1 {
		fmt.Println(len(kd.m))
		t.Fail()
	}
}

func TestLookupData(t *testing.T) {
	kd := NewKademliaStruct()
	token := []byte("AA")
	fmt.Println(token)
	hash := kd.Store(token)
	response := kd.LookupData(hash)
	if response == nil {
		t.Fail()
	}

	fmt.Println(token)
	response = kd.LookupData(*NewRandomKademliaID())

	if response != nil {
		t.Fail()
	}
}
