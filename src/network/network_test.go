package network

import (
	"fmt"
	"testing"

	"github.com/D7024E-Distributed-Systems/Kademlia/src/d7024e"
)

func TestNetworkStruct(t *testing.T) {
	target := d7024e.NewRandomKademliaID()
	contact := d7024e.NewContact(target, "127.0.0.1:3000")

	network := NewNetwork(&contact)

	fmt.Println(network.CurrentNode.ID)
}
