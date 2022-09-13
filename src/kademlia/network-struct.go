package kademlia

type Network struct {
	Kademlia     *Kademlia
	CurrentNode  *Contact
	RoutingTable *RoutingTable
}

func NewNetwork(node *Contact) *Network {
	return &Network{NewKademliaStruct(), node, NewRoutingTable(*node)}
}
