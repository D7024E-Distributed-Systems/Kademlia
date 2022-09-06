package d7024e

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestInsertData(t *testing.T) {
	kd := NewKademliaStruct()
	if len(kd.m) != 0 {
		fmt.Println(len(kd.m))
		t.Fail()
	}

	kd.Store(make([]byte, 255))

	if len(kd.m) != 1 {
		fmt.Println(len(kd.m))
		t.Fail()
	}
}

func TestLookupData(t *testing.T) {
	kd := NewKademliaStruct()
	token := make([]byte, 255)
	rand.Read(token)
	fmt.Println(token)
	hash := kd.Store(token)
	response := kd.LookupData(hash)
	if response == nil {
		t.Fail()
	}

	rand.Read(token)
	fmt.Println(token)
	response = kd.LookupData(Hash(token))

	if response != nil {
		t.Fail()
	}
}
