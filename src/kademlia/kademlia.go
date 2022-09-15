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

func (kademlia *Kademlia) LookupContact(target *KademliaID) *Contact {
	if target.Equals(kademlia.Network.RoutingTable.me.ID) {
		return &kademlia.Network.RoutingTable.me
	}
	contacts := kademlia.Network.RoutingTable.FindClosestContacts(target, alpha)
	fmt.Println("Find closest contacts: ", contacts)
	for _, contact := range contacts {
		if target.Equals(contact.ID) {
			return &contact
		}
	}
	return kademlia.lookupContactHelper(target, contacts)
}

func (kademlia *Kademlia) lookupContactHelper(target *KademliaID, previousContacts []Contact) *Contact {
	routingTable := NewRoutingTable(*kademlia.Network.CurrentNode)
	fmt.Println("previousContacts: ", previousContacts)
	fmt.Println("previousContacts len: ", len(previousContacts))
	for _, contact := range previousContacts {
		// routingTable.AddContact(contact)
		fmt.Println("Sending find contact message to contact:", contact)
		fetchedContacts := kademlia.Network.SendFindContactMessage(&contact, target)
		fmt.Println("Found contacts from fetchedContacts", fetchedContacts)
		for _, tempContact := range fetchedContacts {
			if target.Equals(tempContact.ID) {
				return &tempContact
			}
			routingTable.AddContact(tempContact)
		}
	}
	closestContacts := routingTable.FindClosestContacts(target, alpha)
	fmt.Println("Closest contacts: ", closestContacts)
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
		return nil
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
