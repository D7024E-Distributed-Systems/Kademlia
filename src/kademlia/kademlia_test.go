package kademlia

import (
	"fmt"
	"testing"
	"time"
)

func TestLessThan(t *testing.T) {
	nodeID := NewKademliaID("A")
	nodeID2 := NewKademliaID("B")
	res := nodeID.Less(nodeID2)
	if !res {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {
	nodeID := NewKademliaID("A")
	nodeID2 := NewKademliaID("A")
	res := nodeID.Equals(nodeID2)
	if !res {
		t.Fail()
	}
}

func TestCalcDistance(t *testing.T) {
	nodeID := NewKademliaID("K")
	nodeID2 := NewKademliaID("AAAA")
	res := nodeID.CalcDistance(nodeID2)
	if res.String() != "04055054010955505600030c05000d52070e5e07" {
		t.Fail()
	}
}

func TestInsertData(t *testing.T) {
	kd := NewKademliaStruct()
	if len(kd.M) != 0 {
		fmt.Println(len(kd.M))
		t.Fail()
	}

	kd.Store([]byte("AA"), time.Minute)

	if len(kd.M) != 1 {
		fmt.Println(len(kd.M))
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

	if len(kd.M) != 0 {
		t.Fail()
	}
}

func TestContacts(t *testing.T) {
	contact := NewContact(NewKademliaID("A"), "localhost")
	contact2 := NewContact(NewKademliaID("B"), "localhost")
	hash := NewKademliaID("B")
	hash2 := NewKademliaID("B")
	contact.CalcDistance(hash)
	contact2.CalcDistance(hash2)

	contact.Less(&contact2)

	if contact.String() != "contact(\"3664636434636532336438386532656539353638\", \"localhost\")" {
		t.Fail()
	}
}

func TestContactCandidates(t *testing.T) {
	contact := NewContact(NewKademliaID("A"), "localhost")
	contact2 := NewContact(NewKademliaID("B"), "localhost")
	contactCan := ContactCandidates{}

	contactCan.Append([]Contact{contact, contact2})
	con := contactCan.GetContacts(1)
	if con[0] != contact {
		t.Fail()
	}
	len := contactCan.Len()
	if len != 2 {
		t.Fail()
	}
	contactCan.Swap(0, 1)
}
