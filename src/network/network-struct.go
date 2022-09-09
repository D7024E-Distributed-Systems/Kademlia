package network

import (
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/routing"
)

type Network struct {
	Kademlia     *Kademlia
	CurrentNode  *Contact
	RoutingTable *RoutingTable
}

func NewNetwork(node *Contact) *Network {
	return &Network{NewKademliaStruct(), node, NewRoutingTable(*node)}
}
