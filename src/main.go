package main

import (
	"time"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/d7024e"
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/network"
)

func main() {
	target := NewRandomKademliaID()
	contact := NewContact(target, "127.0.0.1:3000")
	contact.CalcDistance(target)
	currentContact := NewContact(NewRandomKademliaID(), "127.0.0.1:3000")
	network := NewNetwork(&currentContact)
	ping := network.SendPingMessage(&contact)
	if ping {
		// TODO: Send find node to server and start listening on random port
	} else {
		go Listen("127.0.0.1", 3000)
		for {
			time.Sleep(2 * time.Second)
		}
	}
}
