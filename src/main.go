package main

import (
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/d7024e"
)

func main() {
	// go Listen("127.0.0.1", 3000)
	target := NewRandomKademliaID()
	contact := NewContact(target, "127.0.0.1:3000")
	contact.CalcDistance(target)
	currentContact := NewContact(NewRandomKademliaID(), "127.0.0.1:3000")
	network := NewNetwork(&currentContact)
	// time.Sleep(2 * time.Second)
	network.SendPingMessage(&contact)
}
