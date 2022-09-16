package kademlia

import (
	"fmt"
	"testing"
	"time"
)

// To run tests:
// go test -coverprofile cover.out =./... ./...
// To see coverage:
// go tool cover -html=cover.out

func TestNetworkStruct(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:3000")

	network := NewNetwork(&contact)

	fmt.Println(network.CurrentNode.ID)
}

func TestPingNode(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8000")
	network := NewNetwork(&contact)

	go network.Listen("127.0.0.1", 8000)

	go network.SendPingMessage(&contact)

	time.Sleep(1 * time.Millisecond)

	return
}

func TestFindNode(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8000")
	network := NewNetwork(&contact)

	go network.Listen("127.0.0.1", 8000)

	go network.SendFindContactMessage(&contact, nodeID)

	time.Sleep(1 * time.Millisecond)

	return
}

func TestStoreAndFind(t *testing.T) {
	nodeID := NewRandomKademliaID()
	nodeID2 := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8001")
	contact2 := NewContact(nodeID2, "127.0.0.1:8002")
	network := NewNetwork(&contact)
	network2 := NewNetwork(&contact2)

	go network.Listen("127.0.0.1", 8001)
	go network2.Listen("127.0.0.1", 8002)

	time.Sleep(1 * time.Millisecond)

	network2.SendStoreMessage([]byte("String"), 5*time.Second, &contact)

	time.Sleep(1 * time.Millisecond)

	network.Kademlia.DeleteOldData()

	hash := NewKademliaID("String")
	res := network2.SendFindDataMessage(hash, &contact)
	if res != "String" {
		t.Fail()
	}
	time.Sleep(6 * time.Second)
	network.Kademlia.DeleteOldData()
	res2 := network2.SendFindDataMessage(hash, &contact)

	if res2 == "String" {
		t.Fail()
	}

	return
}

func TestFind(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8000")
	network := NewNetwork(&contact)

	go network.Listen("127.0.0.1", 8000)

	time.Sleep(1 * time.Millisecond)

	go network.SendFindDataMessage(nodeID, &contact)

	time.Sleep(1 * time.Millisecond)

	return
}

func TestStoreAndFindAndRefresh(t *testing.T) {
	nodeID := NewRandomKademliaID()
	nodeID2 := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8001")
	contact2 := NewContact(nodeID2, "127.0.0.1:8002")
	network := NewNetwork(&contact)
	network2 := NewNetwork(&contact2)

	go network.Listen("127.0.0.1", 8001)
	go network2.Listen("127.0.0.1", 8002)

	time.Sleep(1 * time.Millisecond)

	network2.SendStoreMessage([]byte("String"), 5*time.Second, &contact)

	time.Sleep(1 * time.Millisecond)

	network.Kademlia.DeleteOldData()

	hash := NewKademliaID("String")
	res := network2.SendFindDataMessage(hash, &contact)
	if res != "String" {
		t.Fail()
	}

	network2.SendRefreshMessage(hash, &contact)
	time.Sleep(4 * time.Second)
	network.Kademlia.DeleteOldData()
	res2 := network2.SendFindDataMessage(hash, &contact)

	if res2 != "String" {
		t.Fail()
	}

	return
}
