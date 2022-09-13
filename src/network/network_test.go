package network

import (
	"fmt"
	"testing"
	"time"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
)

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
	contact := NewContact(nodeID, "127.0.0.1:8000")
	network := NewNetwork(&contact)

	go network.Listen("127.0.0.1", 8000)

	go network.SendStoreMessage([]byte("String"), &contact)

	time.Sleep(1 * time.Millisecond)

	hash := NewKademliaID("String")
	go network.SendFindDataMessage(hash, &contact)

	time.Sleep(1 * time.Millisecond)

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
