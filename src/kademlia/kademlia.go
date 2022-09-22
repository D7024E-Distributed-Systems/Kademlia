package kademlia

import (
	"sync"
	"time"
)

const alpha = 3
const defaultTTL = 10

type Kademlia struct {
	storeValues  map[KademliaID]*Value
	storeMutex   sync.Mutex
	Network      *Network
	KnownHolders map[Contact]KademliaID
	holderMutex  sync.Mutex
}

type Value struct {
	Data               []byte
	timeSinceRepublish int
	TTL                time.Duration
	DeadAt             time.Time
}

func NewKademliaStruct(network *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.storeValues = make(map[KademliaID]*Value)
	kademlia.Network = network
	kademlia.KnownHolders = make(map[Contact]KademliaID)
	kademlia.storeMutex = sync.Mutex{}
	kademlia.holderMutex = sync.Mutex{}
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
	kademlia.storeMutex.Lock()
	defer kademlia.storeMutex.Unlock()
	value, exists := kademlia.storeValues[hash]
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
	}
	return nil, Contact{}
	
func (kademlia *Kademlia) StoreRPC(data []byte) []*KademliaID {
	target := HashDataReturnKademliaID(string(data))
	closest := kademlia.LookupContact(target)
	var storedNodes []*KademliaID
	for _, contact := range closest {
		if contact.ID.Equals(kademlia.Network.RoutingTable.me.ID) {
			kademlia.Store(data, defaultTTL)
			continue
		}
		res := kademlia.Network.SendStoreMessage(data, defaultTTL, &contact, kademlia)
		if res {
			storedNodes = append(storedNodes, contact.ID)
		}
	}
	return storedNodes
}

// Stores data in this node, returns hash of object
func (kademlia *Kademlia) Store(data []byte, ttl time.Duration) (KademliaID, time.Time) {
	hash := NewKademliaID(string(data))
	file := Value{data, 0, ttl, time.Now().Add(ttl)}
	// Mutex lock
	kademlia.storeMutex.Lock()
	defer kademlia.storeMutex.Unlock()
	kademlia.storeValues[*hash] = &file
	return *hash, file.DeadAt
}

func (kademlia *Kademlia) DeleteOldData() {
	kademlia.storeMutex.Lock()
	defer kademlia.storeMutex.Unlock()
	for hash, value := range kademlia.storeValues {
		if time.Now().After(value.DeadAt) {
			delete(kademlia.storeValues, hash)
		}
	}
}

func (kademlia *Kademlia) RefreshTTL(hash KademliaID) {
	kademlia.storeMutex.Lock()
	defer kademlia.storeMutex.Unlock()
	value, exists := kademlia.storeValues[hash]
	if exists {
		value.DeadAt = time.Now().Add(value.TTL)
	}
}

func (kademlia *Kademlia) AddToKnown(contact *Contact, hash *KademliaID) {
	kademlia.holderMutex.Lock()
	defer kademlia.holderMutex.Unlock()
	kademlia.KnownHolders[*contact] = *hash
}

func (kademlia *Kademlia) RemoveFromKnown(value string) bool {
	kademlia.holderMutex.Lock()
	defer kademlia.holderMutex.Unlock()
	kademliaID := ToKademliaID(value)
	for contact, data := range kademlia.KnownHolders {
		if data == *kademliaID {
			delete(kademlia.KnownHolders, contact)
			return true
		}
	}
	return false
}
