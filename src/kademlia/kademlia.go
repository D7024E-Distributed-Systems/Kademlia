package kademlia

import (
	"fmt"
	"time"
)

const alpha = 3

type Kademlia struct {
	M            map[KademliaID]*Value
	KnownHolders map[KademliaID]Contact
}

type Value struct {
	data               []byte
	timeSinceRepublish int
	TTL                time.Duration
	deadAt             time.Time
}

func NewKademliaStruct() *Kademlia {
	kademlia := &Kademlia{}
	kademlia.M = make(map[KademliaID]*Value)
	kademlia.KnownHolders = make(map[KademliaID]Contact)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
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
	hash := NewKademliaID(string(data))
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
	kademlia.KnownHolders[*hash] = *contact
	fmt.Println("KNOWN ARE:", kademlia.KnownHolders)
}
