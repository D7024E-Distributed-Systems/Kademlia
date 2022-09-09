package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
)

const alpha = 3

type Kademlia struct {
	m map[string]Value
}

type Value struct {
	data               []byte
	timeSinceRepublish int
}

func NewKademliaStruct() *Kademlia {
	kademlia := &Kademlia{}
	kademlia.m = make(map[string]Value)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

// Checks if data is stored in this node, returns data if found
func (kademlia *Kademlia) LookupData(hash string) []byte {
	value, exists := kademlia.m[hash]
	if exists {
		return value.data
	}
	return nil
}

// Stores data in this node, returns hash of object
func (kademlia *Kademlia) Store(data []byte) string {
	hash := Hash(data)
	file := Value{data, 0}
	kademlia.m[hash] = file
	return hash
}

// Hashes a given byte splice
func Hash(data []byte) string {
	hashbytes := sha1.Sum(data)
	hash := hex.EncodeToString(hashbytes[0:IDLength])
	return hash
}
