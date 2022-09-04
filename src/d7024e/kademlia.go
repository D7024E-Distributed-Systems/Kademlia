package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	m map[string]Value
}

type Value struct {
	data               []byte
	timeSinceRepublish int
}

func newKademliaStruct() *Kademlia {
	kademlia := &Kademlia{}
	kademlia.m = make(map[string]Value)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	value, exists := kademlia.m[hash]
	if exists {
		return value.data
	}
	return nil
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	hash := Hash(data)
	file := Value{data, 0}
	kademlia.m[hash] = file
	// TODO
}

func Hash(data []byte) string {
	hashbytes := sha1.Sum(data)
	hash := hex.EncodeToString(hashbytes[0:IDLength])
	return hash
}
