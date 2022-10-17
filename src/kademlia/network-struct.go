package kademlia

type Network struct {
	CurrentNode  *Contact
	RoutingTable *RoutingTable
}

/*
Returns a new network instance
  - node *Contact, the current node to be added as the routing table current node
*/
func NewNetwork(node *Contact) *Network {
	return &Network{node, NewRoutingTable(*node)}
}
