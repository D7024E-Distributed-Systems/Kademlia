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

func (kademlia *Kademlia) LookupContact(target *KademliaID) ContactCandidates {
	contacts := kademlia.Network.RoutingTable.FindClosestContacts(target, BucketSize)
	contactCandidates := kademlia.lookupContactHelper(target, contacts)
	allContacts := contactCandidates.contacts
	if target.Equals(kademlia.Network.RoutingTable.me.ID) || 20 > len(allContacts) {
		contact := kademlia.Network.RoutingTable.me
		contact.CalcDistance(target)
		allContacts = append([]Contact{contact}, allContacts...)
		contactCandidates := ContactCandidates{allContacts}
		contactCandidates.Sort()
		return contactCandidates
	}
	return contactCandidates
}

func (kademlia *Kademlia) lookupContactHelper(target *KademliaID, previousContacts []Contact) ContactCandidates {
	routingTable := NewRoutingTable(*kademlia.Network.CurrentNode)
	routingTableLock := sync.Mutex{}
	var wg sync.WaitGroup
	length := min(alpha, len(previousContacts))
	wg.Add(length)
	for i := 0; i < length; i++ {
		contact := previousContacts[i]
		go func(contact Contact) {
			defer wg.Done()
			fetchedContacts := kademlia.Network.SendFindContactMessage(&contact, target)
			routingTableLock.Lock()
			defer routingTableLock.Unlock()
			for _, tempContact := range fetchedContacts {
				routingTable.AddContact(tempContact)
			}
		}(contact)
	}
	wg.Wait()
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
		return ContactCandidates{closestContacts}
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
	candidates := kademlia.LookupContact(hash).contacts
	for len(candidates) > 0 {
		length := min(alpha, len(candidates))
		var wg sync.WaitGroup
		wg.Add(length)
		var resString *string = nil
		var resCandidate Contact = Contact{}
		for i := 0; i < length; i++ {
			go func(candidate Contact) {
				defer wg.Done()
				res := kademlia.Network.SendFindDataMessage(hash, &candidate)
				if !(res == "Error: Invalid contact information" || res == "ERROR" || res == "") {
					// no need for mutex lock since if we get here they will all return the same value
					resString = &res
					resCandidate = candidate
				}
			}(candidates[0])
			candidates = candidates[1:]
		}
		wg.Wait()
		if resString != nil {
			return resString, resCandidate
		}
	}
	return nil, Contact{}
}

// Sends store RPCs to nodes that should store the data
func (kademlia *Kademlia) StoreValue(data []byte, ttl time.Duration) ([]*KademliaID, string) {
	target := NewKademliaID(string(data))
	closest := kademlia.LookupContact(target)
	var storedNodes []*KademliaID
	storedNodesMutex := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(closest.contacts))
	for _, contact := range closest.contacts {
		if contact.ID.Equals(kademlia.Network.RoutingTable.me.ID) {
			kademlia.Store(data, ttl)
			storedNodesMutex.Lock()
			storedNodes = append(storedNodes, contact.ID)
			storedNodesMutex.Unlock()
			wg.Done()
			continue
		}
		go func(contact Contact) {
			defer wg.Done()
			res := kademlia.Network.SendStoreMessage(data, ttl, &contact, kademlia)
			if res {
				storedNodesMutex.Lock()
				storedNodes = append(storedNodes, contact.ID)
				storedNodesMutex.Unlock()
			}
		}(contact)
	}
	wg.Wait()
	return storedNodes, target.String()
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

func (kademlia *Kademlia) DeleteOldDataLoop() {
	for {
		kademlia.DeleteOldData()
		time.Sleep(1 * time.Second)
	}
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

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
