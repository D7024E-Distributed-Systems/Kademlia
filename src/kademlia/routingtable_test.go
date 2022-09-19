package kademlia

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	closestIDs := []string{"6234626662306365376530653666313663326163",
		"6234333338303562326338643739326337356535", "6231323231643738646463656364303938316238", "6262393837313061343962386332613962393061",
		"6566323865383337366164633634613632636130", "3533613163656335373035313638393064326631",
	}

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].ID.String(), closestIDs[i])
		if contacts[i].ID.String() != closestIDs[i] {
			t.FailNow()
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
