package kademlia

import (
	"fmt"
	"time"
)

const alpha = 3

type Kademlia struct {
	m            map[KademliaID]Value
	KnownHolders map[Contact]KademliaID
}

type Value struct {
	data               []byte
	timeSinceRepublish int
	TTL                time.Duration
	deadAt             time.Time
}

func NewKademliaStruct() *Kademlia {
	kademlia := &Kademlia{}
	kademlia.m = make(map[KademliaID]Value)
	kademlia.KnownHolders = make(map[Contact]KademliaID)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

// Checks if data is stored in this node, returns data if found
func (kademlia *Kademlia) LookupData(hash KademliaID) []byte {
	value, exists := kademlia.m[hash]
	if exists {
		// value.deadAt = time.Now().Add(value.TTL)
		return value.data
	}
	return nil
}

// Stores data in this node, returns hash of object
func (kademlia *Kademlia) Store(data []byte, ttl time.Duration) KademliaID {
	hash := NewKademliaID(string(data))
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

func (kademlia *Kademlia) RefreshTTL(hash KademliaID) {
	value, exists := kademlia.m[hash]
	if exists {
		fmt.Println("REFRESHED!")
		value.deadAt = time.Now().Add(value.TTL)
	}
}

func (kademlia *Kademlia) AddToKnown(contact *Contact, hash *KademliaID) {
	kademlia.KnownHolders[*contact] = *hash
}
