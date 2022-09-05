package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/D7024E-Distributed-Systems/Kademlia/src/d7024e"
)

func TestNetworkStruct(t *testing.T) {
	nodeID := d7024e.NewRandomKademliaID()
	contact := d7024e.NewContact(nodeID, "127.0.0.1:3000")

	network := NewNetwork(&contact)

	fmt.Println(network.CurrentNode.ID)
}

func TestPingNode(t *testing.T) {
	nodeID := d7024e.NewRandomKademliaID()
	contact := d7024e.NewContact(nodeID, "127.0.0.1:8000")
	network := NewNetwork(&contact)

	go network.Listen("127.0.0.1", 8000)

	go network.SendPingMessage(&contact)

	time.Sleep(1 * time.Millisecond)

	return
}
