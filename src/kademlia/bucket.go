package kademlia

import (
	"container/list"
	"sync"
)

/*
Bucket struct
  - list *list.List, a list of contacts
*/
type bucket struct {
	list *list.List
}

/*
Returns a new bucket struct
*/
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

/*
AddContact adds a contact to the front of the bucket if not already present, else it moves it to the front
  - contact Contact, the contact to be added
  - self Contact, the current node, used to remove stale contacts
  - mutex *sync.Mutex, mutex lock for thread safety
*/
func (bucket *bucket) AddContact(contact Contact, self Contact, mutex *sync.Mutex) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < BucketSize {
			bucket.list.PushFront(contact)
		} else {
			go func() {
				if bucket.pingAlphaNodesAndRemove(self, mutex) {
					mutex.Lock()
					bucket.list.PushFront(contact)
					mutex.Unlock()
				}
			}()
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

/*
Returns an array of contacts where the distance has already been calculated
  - target *KademliaID, the target to calculate the distance to
  - []Contact, the array of contacts with distance calculated
*/
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

/*
Returns the length of a bucket
*/
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}

/*
Removes contacts, starting from the back of a bucket, and goes alpha deep into the bucket.
It will remove the first contact which does not respond
  - self Contact, creates a new network with itself as base
  - mutex *sync.Mutex, mutex lock for thread safety
*/
func (bucket *bucket) pingAlphaNodesAndRemove(self Contact, mutex *sync.Mutex) bool {
	network := NewNetwork(&self)
	i := 0
	length := min(alpha, bucket.Len())
	for e := bucket.list.Back(); e != nil; e = e.Prev() {
		if i >= length {
			break
		}
		contact := e.Value.(Contact)
		if !network.SendPingMessage(&contact) {
			mutex.Lock()
			bucket.list.Remove(e)
			mutex.Unlock()
			return true
		}
		i++
	}
	return false
}
