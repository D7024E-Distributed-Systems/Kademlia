package d7024e

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestInsertData(t *testing.T) {
	kd := newKademliaStruct()
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
	kd := newKademliaStruct()
	token := make([]byte, 255)
	rand.Read(token)
	fmt.Println(token)
	kd.Store(token)
	response := kd.LookupData(Hash(token))
	if response == nil {
		t.Fail()
	}

	token2 := make([]byte, 255)
	rand.Read(token2)
	fmt.Println(token2)
	response = kd.LookupData(Hash(token2))

	if response != nil {
		t.Fail()
	}
}
