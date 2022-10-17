package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// the static number of bytes in a KademliaID
const IDLength = 20

// type definition of a KademliaID
type KademliaID [IDLength]byte

/*
Makes a new KademliaID from a string.
  - data string, the value that should be hashed to a KademliaID
*/
func NewKademliaID(data string) *KademliaID {
	hashBytes := sha1.Sum([]byte(data))
	hash := hex.EncodeToString(hashBytes[0:IDLength])

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = hash[i]
	}

	return &newKademliaID
}

/*
Casts a string to a KademliaID. Note that this is not the same as NewKademliaID().
  - bar string, the string to cast
*/
func ToKademliaID(bar string) *KademliaID {
	if len(bar) < 40 {
		return nil
	}
	res, err := hex.DecodeString(bar)
	if err != nil {
		fmt.Println("FAILED TO DECODE KADEMLIA ID", err)
		return nil
	} else {
		return (*KademliaID)(res)
	}
}

/*
Makes a random new KademliaID
*/
func NewRandomKademliaID() *KademliaID {
	rand.Seed(time.Now().UnixNano())
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

/*
Returns true if kademliaID < otherKademliaID (bitwise)
  - otherKademliaID *KademliaID, the KademliaID to be compared to
*/
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals retuns trrue if kademliaID == otherKademliaID (bitwise)
/*
Returns trrue if kademliaID == otherKademliaID (bitwise)
  - otherKademliaID *KademliaID, the KademliaID to be compared to
*/
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

/*
Returns the distance between two KademliaIDs. Uses bitwise XOR.
  - target *KademliaID, the target to calculate the distance to
*/
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

/*
Returns a string representation of a KademliaID
*/
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
