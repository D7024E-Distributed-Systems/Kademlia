package kademlia

import "sync"

const BucketSize = 20

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	me          Contact
	buckets     [IDLength * 8]*bucket
	bucketMutex sync.Mutex
}

/*
Returns a new instance of a RoutingTable

  - me Contact, sets the current contact to ourself
*/
func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	routingTable.bucketMutex = sync.Mutex{}
	return routingTable
}

/*
Adds a new contact to the correct Bucket from our routing table
  - contact Contact, add this contact to the bucket if the bucket is not full
*/
func (routingTable *RoutingTable) AddContact(contact Contact) {
	routingTable.bucketMutex.Lock()
	defer routingTable.bucketMutex.Unlock()
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContact(contact, routingTable.me, &routingTable.bucketMutex)
}

/*
Finds the "count" closest Contact to the target from the routing table
  - target *KademliaID, the id to calculate the closest distance to
  - count int, the number of contacts to return
*/
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	var candidates ContactCandidates
	routingTable.bucketMutex.Lock()
	defer routingTable.bucketMutex.Unlock()
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

/*
Get the correct Bucket index for the KademliaID
  - id *KademliaID, the id of a node to find the bucket index for
*/
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}
