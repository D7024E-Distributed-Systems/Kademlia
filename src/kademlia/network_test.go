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
	kademlia := NewKademliaStruct(network)

	go network.Listen("127.0.0.1", 8000, kademlia)

	go network.SendPingMessage(&contact)

	time.Sleep(1 * time.Millisecond)

	return
}

func TestPingNode2(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.:8000")
	network := NewNetwork(&contact)

	if network.SendPingMessage(&contact) {
		t.Fail()
	}
}

func TestFindNode(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8001")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	go network.Listen("127.0.0.1", 8001, kademlia)

	go network.SendFindContactMessage(&contact, nodeID)

	time.Sleep(1 * time.Millisecond)

	return
}

func TestFindNode2(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.:8000")
	network := NewNetwork(&contact)

	res := network.SendFindContactMessage(&contact, NewRandomKademliaID())
	if res != nil {
		t.Fail()
	}
}

func TestFindDataMessage(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.:8000")
	network := NewNetwork(&contact)

	res := network.SendFindDataMessage(NewRandomKademliaID(), &contact)
	if res != "ERROR" {
		t.Fail()
	}
}

func TestSendStoreMessage(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.:8000")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	res := network.SendStoreMessage([]byte("test"), time.Duration(1*time.Second), &contact, kademlia)
	if res != false {
		t.Fail()
	}
}

func TestSendRefreshMessage(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.:8000")
	network := NewNetwork(&contact)

	res := network.SendRefreshMessage(NewRandomKademliaID(), &contact)
	if res != false {
		t.Fail()
	}
}

func TestStoreAndFind(t *testing.T) {
	nodeID := NewRandomKademliaID()
	time.Sleep(1 * time.Millisecond)
	nodeID2 := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8002")
	contact2 := NewContact(nodeID2, "127.0.0.1:8003")
	network := NewNetwork(&contact)
	network2 := NewNetwork(&contact2)
	kademlia := NewKademliaStruct(network)
	kademlia2 := NewKademliaStruct(network2)

	go network.Listen("127.0.0.1", 8002, kademlia)
	go network2.Listen("127.0.0.1", 8003, kademlia2)

	time.Sleep(1 * time.Millisecond)
	fmt.Println(network2.SendStoreMessage([]byte("String"), 5*time.Second, &contact, kademlia2))

	time.Sleep(1 * time.Millisecond)

	kademlia.DeleteOldData()

	hash := NewKademliaID("String")
	res := network2.SendFindDataMessage(hash, &contact)
	if res != "String" {
		fmt.Println("Res is", res)
		t.Fail()
	}
	time.Sleep(6 * time.Second)
	kademlia.DeleteOldData()
	res2 := network2.SendFindDataMessage(hash, &contact)

	if res2 == "String" {
		fmt.Println("Res2 is", res)
		t.Fail()
	}

	return
}

func TestFind(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8004")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	go network.Listen("127.0.0.1", 8004, kademlia)

	time.Sleep(1 * time.Millisecond)

	go network.SendFindDataMessage(nodeID, &contact)

	time.Sleep(1 * time.Millisecond)

	return
}

func TestStoreAndFindAndRefresh(t *testing.T) {
	t.Parallel()
	nodeID := NewRandomKademliaID()
	time.Sleep(100 * time.Millisecond)
	nodeID2 := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8005")
	contact2 := NewContact(nodeID2, "127.0.0.1:8006")
	network := NewNetwork(&contact)
	network2 := NewNetwork(&contact2)
	kademlia := NewKademliaStruct(network)
	kademlia2 := NewKademliaStruct(network2)

	go network.Listen("127.0.0.1", 8005, kademlia)
	go network2.Listen("127.0.0.1", 8006, kademlia2)

	time.Sleep(1 * time.Millisecond)

	kademlia2.Network.SendStoreMessage([]byte("String"), 5*time.Second, &contact, kademlia2)

	time.Sleep(1 * time.Millisecond)

	kademlia.DeleteOldData()

	hash := NewKademliaID("String")
	res := kademlia2.Network.SendFindDataMessage(hash, &contact)
	if res != "String" {
		t.Fail()
	}

	kademlia2.Network.SendRefreshMessage(hash, &contact)
	time.Sleep(4 * time.Second)
	kademlia.DeleteOldData()
	res2 := network2.SendFindDataMessage(hash, &contact)

	if res2 != "String" {
		t.Fail()
	}

	return
}

func TestStoreAndFindAndRefreshLoop(t *testing.T) {
	t.Parallel()
	nodeID := NewRandomKademliaID()
	time.Sleep(100 * time.Millisecond)
	nodeID2 := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8005")
	contact2 := NewContact(nodeID2, "127.0.0.1:8006")
	network := NewNetwork(&contact)
	network2 := NewNetwork(&contact2)
	kademlia := NewKademliaStruct(network)
	kademlia2 := NewKademliaStruct(network2)

	go network.Listen("127.0.0.1", 8005, kademlia)
	go network2.Listen("127.0.0.1", 8006, kademlia2)

	time.Sleep(1 * time.Millisecond)

	kademlia2.Network.SendStoreMessage([]byte("String"), 5*time.Second, &contact, kademlia2)

	time.Sleep(1 * time.Millisecond)

	go kademlia2.Network.RefreshLoop(kademlia2)

	kademlia.DeleteOldData()

	hash := NewKademliaID("String")
	res := kademlia2.Network.SendFindDataMessage(hash, &contact)
	if res != "String" {
		t.Fail()
	}

	kademlia2.Network.SendRefreshMessage(hash, &contact)
	time.Sleep(4 * time.Second)
	kademlia.DeleteOldData()
	res2 := network2.SendFindDataMessage(hash, &contact)

	if res2 != "String" {
		t.Fail()
	}

	return
}

func TestListenFailure(t *testing.T) {
	defer func() { recover() }()

	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:3000")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	network.Listen("asdasd", 3, kademlia)

	t.Errorf("did not panic")
}

func TestNonexistentRPC(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:3000")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	res := getResponseMessage([]byte("NONE"), kademlia)

	if string(res) != "Error: Invalid RPC protocol" {
		t.Fail()
	}
}

func TestRefreshResponseFailure(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:3000")
	network := NewNetwork(&contact)

	handleRefreshResponse([]byte("Error"), network)
}

func TestHandleStoreResponseFailure(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:3000")
	network := NewNetwork(&contact)

	handleStoreResponse([]byte("Error"), network)
}

func TestMarshalCurrentNodeFailure(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "")
	network := NewNetwork(&contact)

	network.marshalCurrentNode()
}
