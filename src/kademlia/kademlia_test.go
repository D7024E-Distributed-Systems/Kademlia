package kademlia

import (
	"fmt"
	"testing"
	"time"
)

func TestLessThan(t *testing.T) {
	nodeID := ToKademliaID("A000000000000000000000000000000000000000")
	nodeID2 := ToKademliaID("B000000000000000000000000000000000000000")
	res := nodeID.Less(nodeID2)
	if !res {
		t.Fail()
	}
}

func TestLessThanEqual(t *testing.T) {
	nodeID := ToKademliaID("A000000000000000000000000000000000000000")
	nodeID2 := ToKademliaID("A000000000000000000000000000000000000000")
	res := nodeID2.Less(nodeID)
	if res {
		t.Fail()
	}
}

func TestKademliaId(t *testing.T) {
	node := ToKademliaID("A000000000000000000000000000000000000000")
	if node.String() != "a000000000000000000000000000000000000000" {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {
	nodeID := ToKademliaID("A000000000000000000000000000000000000000")
	nodeID2 := ToKademliaID("A000000000000000000000000000000000000000")
	res := nodeID.Equals(nodeID2)
	if !res {
		t.Fail()
	}
}

func TestCalcDistance(t *testing.T) {
	nodeID := ToKademliaID("A0A0A00000000000000000000000000000000000")
	nodeID2 := ToKademliaID("AAAAAAAAAAAAAAAAA00000000000000000000000")
	res := nodeID.CalcDistance(nodeID2)
	fmt.Println(res)
	if res.String() != "0a0a0aaaaaaaaaaaa00000000000000000000000" {
		t.Fail()
	}
}

func TestInsertData(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost")
	kd := NewKademliaStruct(NewNetwork(&contact))
	if len(kd.storeValues) != 0 {
		fmt.Println(len(kd.storeValues))
		t.Fail()
	}

	kd.Store([]byte("AA"), time.Minute)

	if len(kd.storeValues) != 1 {
		fmt.Println(len(kd.storeValues))
		t.Fail()
	}
}

func TestLookupData(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost")
	kd := NewKademliaStruct(NewNetwork(&contact))
	token := []byte("AA")
	fmt.Println(token)
	hash, _ := kd.Store(token, time.Minute)
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

func TestStoreData(t *testing.T) {
	kademliaNodes := returnKademliaNodes()
	res, _ := kademliaNodes[3].StoreValue([]byte("lmao"), time.Minute)
	if len(res) != len(kademliaNodes) {
		fmt.Println("STORED ON:", res, "\nTOTAL NODES:", len(kademliaNodes))
		t.FailNow()
	}
	for i, node := range kademliaNodes {
		if len(node.storeValues) == 0 {
			fmt.Println("Node", i, "has len", len(node.storeValues))
			t.FailNow()
		}
	}
}

func TestDeleteData(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost")
	kd := NewKademliaStruct(NewNetwork(&contact))
	token := []byte("AA")
	fmt.Println(token)
	kd.Store(token, time.Second)

	time.Sleep(1 * time.Second)
	kd.DeleteOldData()

	if len(kd.storeValues) != 0 {
		t.Fail()
	}
}

func TestContacts(t *testing.T) {
	contact := NewContact(ToKademliaID("A000000000000000000000000000000000000000"), "localhost")
	contact2 := NewContact(ToKademliaID("B000000000000000000000000000000000000000"), "localhost")

	if contact.String() != "contact(\"a000000000000000000000000000000000000000\", \"localhost\")" {
		t.Fail()
	}
	if contact2.String() != "contact(\"b000000000000000000000000000000000000000\", \"localhost\")" {

		t.Fail()
	}
}

func TestContactCandidates(t *testing.T) {
	contact := NewContact(ToKademliaID("A000000000000000000000000000000000000000"), "localhost")
	contact2 := NewContact(ToKademliaID("B000000000000000000000000000000000000000"), "localhost")
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

func TestDistanceBucket(t *testing.T) {
	nodeId := ToKademliaID("A000000000000000000000000000000000000000")
	me := NewContact(nodeId, "")
	contact := NewContact(nodeId, "")
	rt := NewRoutingTable(me)
	rt.AddContact(contact)
	fmt.Println(rt.FindClosestContacts(contact.ID, 1)[0].distance.String())
	if rt.FindClosestContacts(contact.ID, 1)[0].distance.String() != "0000000000000000000000000000000000000000" {
		t.Fail()
	}
}

func TestFindContact(t *testing.T) {
	kademliaNodes := returnKademliaNodes()
	res := kademliaNodes[3].LookupContact(kademliaNodes[0].Network.CurrentNode.ID).contacts
	fmt.Println(res, kademliaNodes[0].Network.CurrentNode.ID)
	if res == nil || !res[0].ID.Equals(kademliaNodes[0].Network.CurrentNode.ID) {
		t.Fail()
	}
}

func TestFindContact2(t *testing.T) {
	kademliaNodes := returnKademliaNodes()
	res := kademliaNodes[3].LookupContact(kademliaNodes[1].Network.CurrentNode.ID).contacts
	if res == nil || !res[0].ID.Equals(kademliaNodes[1].Network.CurrentNode.ID) {
		t.Fail()
	}
}

func TestFindContact3(t *testing.T) {
	kademliaNodes := returnKademliaNodes()
	res := kademliaNodes[3].LookupContact(kademliaNodes[2].Network.CurrentNode.ID).contacts
	if res == nil || !res[0].ID.Equals(kademliaNodes[2].Network.CurrentNode.ID) {
		t.Fail()
	}
}

func TestFindContact4(t *testing.T) {
	kademliaNodes := returnKademliaNodes()
	kadId := ToKademliaID(kademliaNodes[3].Network.CurrentNode.ID.String())
	res := kademliaNodes[3].LookupContact(kadId).contacts
	fmt.Println(res, kademliaNodes[3].Network.CurrentNode)
	if res == nil || !res[0].ID.Equals(kademliaNodes[3].Network.CurrentNode.ID) {
		t.Fail()
	}
}

func returnKademliaNodes() []*Kademlia {
	nodeID := ToKademliaID("A000000000000000000000000000000000000000")
	contact := NewContact(nodeID, "127.0.0.1:7000")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)
	nodeID2 := ToKademliaID("B000000000000000000000000000000000000000")
	contact2 := NewContact(nodeID2, "127.0.0.1:7001")
	network2 := NewNetwork(&contact2)
	kademlia2 := NewKademliaStruct(network2)
	nodeID3 := ToKademliaID("C000000000000000000000000000000000000000")
	contact3 := NewContact(nodeID3, "127.0.0.1:7002")
	network3 := NewNetwork(&contact3)
	kademlia3 := NewKademliaStruct(network3)
	nodeID4 := ToKademliaID("D000000000000000000000000000000000000000")
	contact4 := NewContact(nodeID4, "127.0.0.1:7003")
	network4 := NewNetwork(&contact4)
	kademlia4 := NewKademliaStruct(network4)

	go kademlia.Network.Listen("127.0.0.1", 7000, kademlia)
	go kademlia2.Network.Listen("127.0.0.1", 7001, kademlia2)
	go kademlia3.Network.Listen("127.0.0.1", 7002, kademlia3)
	go kademlia3.Network.Listen("127.0.0.1", 7003, kademlia4)
	kademlia2.Network.RoutingTable.AddContact(contact)
	kademlia3.Network.RoutingTable.AddContact(contact2)
	kademlia4.Network.RoutingTable.AddContact(contact3)
	fmt.Println("Contact 1", kademlia.Network.CurrentNode.ID)
	fmt.Println("Contact 2", kademlia2.Network.CurrentNode.ID)
	fmt.Println("Contact 3", kademlia3.Network.CurrentNode.ID)
	fmt.Println("Contact 4", kademlia4.Network.CurrentNode.ID)
	time.Sleep(1 * time.Second)
	kademliaArray := make([]*Kademlia, 4)
	kademliaArray[0] = kademlia
	kademliaArray[1] = kademlia2
	kademliaArray[2] = kademlia3
	kademliaArray[3] = kademlia4
	return kademliaArray
}

func TestBucketLength(t *testing.T) {
	buck := newBucket()
	if buck.Len() != 0 {
		t.Fail()
	}
	buck.AddContact(NewContact(NewRandomKademliaID(), ""))
	if buck.Len() != 1 {
		t.Fail()
	}
	time.Sleep(1 * time.Millisecond)
	buck.AddContact(NewContact(NewRandomKademliaID(), ""))
	if buck.Len() != 2 {
		t.Fail()
	}

}

func TestKademliaIdNotLess(t *testing.T) {
	id := ToKademliaID("A000000000000000000000000000000000000000")
	id2 := ToKademliaID("b000000000000000000000000000000000000000")
	if id2.Less(id) {
		t.Fail()
	}
}

func TestToKademliaIDSuccess(t *testing.T) {
	id := ToKademliaID("A000000000000000000000000000000000000000")
	expected := ToKademliaID(id.String())
	if *id != *expected {
		t.Fail()
	}
}

func TestToKademliaIDFailure(t *testing.T) {
	actual := ToKademliaID("LL")
	if actual != nil {
		t.Fail()
	}
}

func TestRemoveFromKnownSuccess(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost:3000")
	kademlia := NewKademliaStruct(NewNetwork(&contact))
	kademlia.KnownHolders[contact] = *ToKademliaID("B000000000000000000000000000000000000000")
	kademlia.RemoveFromKnown("B000000000000000000000000000000000000000")
	fmt.Println(kademlia.KnownHolders)
	if len(kademlia.KnownHolders) != 0 {
		t.Fail()
	}
}

func TestRemoveFromKnownFailure(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "localhost:3000")
	kademlia := NewKademliaStruct(NewNetwork(&contact))
	kademlia.KnownHolders[contact] = *ToKademliaID("C000000000000000000000000000000000000000")
	kademlia.RemoveFromKnown("B000000000000000000000000000000000000000")
	fmt.Println(kademlia.KnownHolders)
	if len(kademlia.KnownHolders) != 1 {
		t.Fail()
	}
}
