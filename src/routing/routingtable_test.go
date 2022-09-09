package routing

import (
	"fmt"
	"testing"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	closestIDs := []string{"2111111400000000000000000000000000000000",
		"1111111400000000000000000000000000000000", "1111111100000000000000000000000000000000", "1111111200000000000000000000000000000000",
		"1111111300000000000000000000000000000000", "ffffffff00000000000000000000000000000000",
	}

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
		if contacts[i].ID.String() != closestIDs[i] {
			t.Fail()
		}
	}
}

func TestRandomID(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8001"))

	contacts := rt.FindClosestContacts(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 20)

	if len(contacts) != 1 {
		t.Fail()
	}
}

// func testAddContactSameId(t *testing.T) {
// 	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
// 	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
// 	rt.AddContact(NewContact(contact), "test")
// 	rt.AddContact(NewContact(contact), "test")
// }
