package kademlia

import (
	"fmt"
	"sort"
)

/*
Contact struct
  - ID *KademliaID, the ID of the node
  - Address string, the IP address and port of the node
  - distance *KademliaID, the distance to the node from the current node
*/
type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}

/*
Returns a new contact
  - id *KademliaID, the ID of the new node
  - address string, the IP address and port of the new node
*/
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

/*
Calculates the distance from the current node to some target
  - target *KademliaID, the ID of the node to calculate the distance from
*/
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

/*
Returns true if contact.distance < otherContact.distance, false otherwise
  - otherContact *Contact, the contact to be compared to
*/
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

// String returns a simple string representation of a Contact
/*
Returns a string representation of a contact
*/
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

/*
ContactCandidates struct, used as a helper struct for keeping track of multiple contacts
  - contats []contact, multiple contacts
*/
type ContactCandidates struct {
	contacts []Contact
}

/*
Appends an array of contacts to a contactCandidates struct
  - contacts []Contact, the contacts to append
*/
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

/*
Returns a number of contacts equal to a number
  - the number of contacts to return
*/
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}

/*
Sorts the contacts in a contactCandidates based on distance
*/
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

/*
Returns the number of contacts in a contactCandidates struct
*/
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
/*
Swaps the position of two contacts in a contactCandidates struct
  Used by the built-in sort function
  - i int, an index
  - j int, an index
*/
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

/*
Returns true if a contact's distance at an index is smaller than the contact's distance at another index
  - i int, an index
  - j int, an index
*/
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}
