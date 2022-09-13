package kademlia

import (
	"fmt"
	"time"
)

const alpha = 3

type Kademlia struct {
	m       map[KademliaID]Value
	Network *Network
}

type Value struct {
	data               []byte
	timeSinceRepublish int
	TTL                time.Duration
	deadAt             time.Time
}

func NewKademliaStruct(network *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.m = make(map[KademliaID]Value)
	kademlia.Network = network
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) *Contact {
	contacts := kademlia.Network.RoutingTable.FindClosestContacts(target, alpha)
	for _, contact := range contacts {
		if target.Equals(contact.ID) {
			return &contact
		}
	}
	return kademlia.lookupContactHelper(target, contacts)
}

func (kademlia *Kademlia) lookupContactHelper(target *KademliaID, previousContacts []Contact) *Contact {
	routingTable := NewRoutingTable(*kademlia.Network.CurrentNode)
	for _, contact := range previousContacts {
		fetchedContacts := kademlia.Network.SendFindContactMessage(&contact, target)
		for _, tempContact := range fetchedContacts {
			if target.Equals(tempContact.ID) {
				return &tempContact
			}
			routingTable.AddContact(contact)
		}
	}
	closestContacts := routingTable.FindClosestContacts(target, alpha)
	howManyContactsKnown := 0
	for _, contact := range closestContacts {
		for _, prevContact := range previousContacts {
			if contact.ID == prevContact.ID {
				howManyContactsKnown++
				break
			}
		}
	}
	if howManyContactsKnown == len(closestContacts) {
		fmt.Println("Exiting find contact since we have gotten to n = n-1")
		return nil
	} else {
		return kademlia.lookupContactHelper(target, closestContacts)
	}
}

// Checks if data is stored in this node, returns data if found
func (kademlia *Kademlia) LookupData(hash KademliaID) []byte {
	value, exists := kademlia.m[hash]
	if exists {
		value.deadAt = time.Now().Add(value.TTL)
		return value.data
	}
	return nil
}

// Stores data in this node, returns hash of object
func (kademlia *Kademlia) Store(data []byte, ttl time.Duration) KademliaID {
	hash := HashDataReturnKademliaID(string(data))
	file := Value{data, 0, ttl, time.Now().Add(ttl)}
	kademlia.m[*hash] = file
	return *hash
}

func (kademlia *Kademlia) DeleteOldData() {
	for hash, value := range kademlia.m {
		if time.Now().After(value.deadAt) {
			delete(kademlia.m, hash)
		}
	}
}
