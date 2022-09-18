package kademlia

import (
	"fmt"
	"time"
)

const alpha = 3

type Kademlia struct {
	M            map[KademliaID]*Value
	Network      *Network
	KnownHolders map[Contact]KademliaID
}

type Value struct {
	data               []byte
	timeSinceRepublish int
	TTL                time.Duration
	deadAt             time.Time
}

func NewKademliaStruct(network *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.M = make(map[KademliaID]*Value)
	kademlia.Network = network
	kademlia.KnownHolders = make(map[Contact]KademliaID)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	contacts := kademlia.Network.RoutingTable.FindClosestContacts(target, BucketSize)
	allContacts := kademlia.lookupContactHelper(target, contacts)
	if target.Equals(kademlia.Network.RoutingTable.me.ID) {
		contact := kademlia.Network.RoutingTable.me
		return append(allContacts, contact)
	}
	return allContacts
}

func (kademlia *Kademlia) lookupContactHelper(target *KademliaID, previousContacts []Contact) []Contact {
	routingTable := NewRoutingTable(*kademlia.Network.CurrentNode)
	for _, contact := range previousContacts {
		fetchedContacts := kademlia.Network.SendFindContactMessage(&contact, target)
		for _, tempContact := range fetchedContacts {
			routingTable.AddContact(tempContact)
		}
	}
	closestContacts := routingTable.FindClosestContacts(target, BucketSize)
	howManyContactsKnown := 0
	for _, contact := range closestContacts {
		for _, prevContact := range previousContacts {
			if contact.ID.Equals(prevContact.ID) {
				howManyContactsKnown++
				break
			}
		}
	}
	if howManyContactsKnown == len(closestContacts) {
		fmt.Println("Exiting find contact since we have gotten to n = n-1")
		return closestContacts
	} else {
		return kademlia.lookupContactHelper(target, closestContacts)
	}
}

// Checks if data is stored in this node, returns data if found
func (kademlia *Kademlia) LookupData(hash KademliaID) []byte {
	value, exists := kademlia.M[hash]
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
	kademlia.M[*hash] = &file
	return *hash
}

func (kademlia *Kademlia) DeleteOldData() {
	for hash, value := range kademlia.M {
		fmt.Println("DEAD IS", value.deadAt)
		if time.Now().After(value.deadAt) {
			delete(kademlia.M, hash)
		}
	}
}

func (kademlia *Kademlia) RefreshTTL(hash KademliaID) {
	value, exists := kademlia.M[hash]
	if exists {
		value.deadAt = time.Now().Add(value.TTL)
	}
}

func (kademlia *Kademlia) AddToKnown(contact *Contact, hash *KademliaID) {
	kademlia.KnownHolders[*contact] = *hash
	fmt.Println("KNOWN ARE:", kademlia.KnownHolders)
}
