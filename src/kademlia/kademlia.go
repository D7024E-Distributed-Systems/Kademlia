package kademlia

import (
	"fmt"
	"time"
)

const alpha = 3

type Kademlia struct {
	m            map[KademliaID]*Value
	Network      *Network
	KnownHolders map[Contact]KademliaID
}

type Value struct {
	Data               []byte
	timeSinceRepublish int
	TTL                time.Duration
	DeadAt             time.Time
}

func NewKademliaStruct(network *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.m = make(map[KademliaID]*Value)
	kademlia.Network = network
	kademlia.KnownHolders = make(map[Contact]KademliaID)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	contacts := kademlia.Network.RoutingTable.FindClosestContacts(target, BucketSize)
	allContacts := kademlia.lookupContactHelper(target, contacts)
	if target.Equals(kademlia.Network.RoutingTable.me.ID) {
		contact := kademlia.Network.RoutingTable.me
		contact.CalcDistance(kademlia.Network.CurrentNode.ID)
		return append([]Contact{contact}, allContacts...)
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
		return closestContacts
	} else {
		return kademlia.lookupContactHelper(target, closestContacts)
	}
}

// Checks if data is stored in this node, returns data if found
func (kademlia *Kademlia) LookupData(hash KademliaID) []byte {
	value, exists := kademlia.m[hash]
	if exists {
		value.DeadAt = time.Now().Add(value.TTL)
		return value.Data
	}
	return nil
}

func (kademlia *Kademlia) GetValue(hash *KademliaID) (*string, Contact) {
	res := kademlia.LookupData(*hash)
	if res != nil {
		ret := string(res)
		return &ret, *kademlia.Network.CurrentNode
	}
	candidates := kademlia.LookupContact(hash)
	for len(candidates) > 0 {
		for i := 0; i < alpha; i++ {
			res := kademlia.Network.SendFindDataMessage(hash, &candidates[0])
			if res != "" {
				return &res, candidates[0]
			}
			candidates = candidates[1:]
		}
		fmt.Println(len(candidates))
	}
	return nil, Contact{}
}

// Stores data in this node, returns hash of object
func (kademlia *Kademlia) Store(data []byte, ttl time.Duration) (KademliaID, time.Time) {
	hash := HashDataReturnKademliaID(string(data))
	file := Value{data, 0, ttl, time.Now().Add(ttl)}
	kademlia.m[*hash] = &file
	return *hash, file.DeadAt
}

func (kademlia *Kademlia) DeleteOldData() {
	for hash, value := range kademlia.m {
		fmt.Println("DEAD IS", value.DeadAt)
		if time.Now().After(value.DeadAt) {
			delete(kademlia.m, hash)
		}
	}
}

func (kademlia *Kademlia) RefreshTTL(hash KademliaID) {
	value, exists := kademlia.m[hash]
	if exists {
		value.DeadAt = time.Now().Add(value.TTL)
	}
}

func (kademlia *Kademlia) AddToKnown(contact *Contact, hash *KademliaID) {
	kademlia.KnownHolders[*contact] = *hash
}

func (kademlia *Kademlia) RemoveFromKnown(value string) bool {
	kademliaID := ToKademliaID(value)
	for contact, data := range kademlia.KnownHolders {
		if data == kademliaID {
			delete(kademlia.KnownHolders, contact)
			return true
		}
	}
	return false
}
