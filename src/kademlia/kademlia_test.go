package kademlia

import (
	"fmt"
	"testing"
	"time"
)

func TestInsertData(t *testing.T) {
	kd := NewKademliaStruct()
	if len(kd.m) != 0 {
		fmt.Println(len(kd.m))
		t.Fail()
	}

	kd.Store([]byte("AA"), time.Minute)

	if len(kd.m) != 1 {
		fmt.Println(len(kd.m))
		t.Fail()
	}
}

func TestLookupData(t *testing.T) {
	kd := NewKademliaStruct()
	token := []byte("AA")
	fmt.Println(token)
	hash := kd.Store(token, time.Minute)
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

func TestDeleteData(t *testing.T) {
	kd := NewKademliaStruct()
	token := []byte("AA")
	fmt.Println(token)
	kd.Store(token, time.Second)

	time.Sleep(1 * time.Second)
	kd.DeleteOldData()

	if len(kd.m) != 0 {
		t.Fail()
	}
}
