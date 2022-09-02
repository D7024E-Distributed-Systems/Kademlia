package network

import (
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/d7024e"
)

type Network struct {
	currentNode *Contact
}

func NewNetwork(node *Contact) *Network {
	return &Network{node}
}
