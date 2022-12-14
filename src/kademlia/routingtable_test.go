package kademlia

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(ToKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	closestIDs := []string{"2111111400000000000000000000000000000000",
		"1111111400000000000000000000000000000000", "1111111100000000000000000000000000000000", "1111111200000000000000000000000000000000",
		"1111111300000000000000000000000000000000", "ffffffff00000000000000000000000000000000",
	}

	rt.AddContact(NewContact(ToKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(ToKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(ToKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(ToKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(ToKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(ToKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(ToKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].ID.String(), closestIDs[i])
		if contacts[i].ID.String() != closestIDs[i] {
			t.FailNow()
		}
	}
}

func TestRandomID(t *testing.T) {
	rt := NewRoutingTable(NewContact(ToKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8001"))

	contacts := rt.FindClosestContacts(ToKademliaID("FFFFFFFF00000000000000000000000000000000"), 20)

	if len(contacts) != 1 {
		t.Fail()
	}
}
