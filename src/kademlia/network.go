package kademlia

import (
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"
)

var maxBytes int = 4096

/*
Sending a ping RPC message to a contact and returns whether it was successful or not
  - contact *Contact, the contact to send the ping RPC to
*/
func (network *Network) SendPingMessage(contact *Contact) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return false
	}
	defer conn.Close()
	message := network.getPingMessage()
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	network.handlePingResponse(buffer[:n])
	return true
}

/*
Returns a standard ping message
*/
func (network *Network) getPingMessage() []byte {
	startMessage := []byte(newPing().startMessage + ";")
	body := network.marshalCurrentNode()
	return append(startMessage, body...)
}

/*
Handle the ping response and parse to see if an error has occurred or if the response was successful
  - message []byte, the response to parse
*/
func (network *Network) handlePingResponse(message []byte) {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return
	} else {
		var contact *Contact
		json.Unmarshal(message, &contact)
		if VerifyContact(contact, network) {
			network.RoutingTable.AddContact(*contact)
		}
	}
}

/*
Sending a send find contact RPC message to a contact and returns the responded list of contacts
  - contact *Contact, the contact to send the RPC to
  - nodeID *KademliaID, the id to find closest contact
*/
func (network *Network) SendFindContactMessage(contact *Contact, nodeID *KademliaID) []Contact {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return nil
	}
	defer conn.Close()
	message := network.getFindContactMessage(nodeID)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return nil
	}
	return network.handleFindContactResponse(buffer[:n])
}

/*
Returns a standard find contact message to find contacts close to nodeID
  - nodeID *KademliaID, the id to find closest contacts to
*/
func (network *Network) getFindContactMessage(nodeID *KademliaID) []byte {
	body, _ := json.Marshal(nodeID)
	startMessage := []byte(newFindContact().startMessage + ";" + string(body) + ";")
	body2 := network.marshalCurrentNode()
	return append(startMessage, body2...)

}

/*
Handle the find contact response and parse to see if an error has occurred or if the response was successful
  - message []byte, the response to parse

Returns the list of contacts
*/
func (network *Network) handleFindContactResponse(message []byte) []Contact {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return nil
	} else {
		var contacts []Contact
		var verifiedContacts []Contact
		json.Unmarshal(message, &contacts)
		for _, contact := range contacts {
			if VerifyContact(&contact, network) {
				if network.SendPingMessage(&contact) {
					verifiedContacts = append(verifiedContacts, contact)
				}
			}
		}
		return verifiedContacts
	}
}

/*
Sending a send find data RPC message to a contact and returns the data if found otherwise "ERROR"
  - hash *KademliaID, the hash to find data for
  - contact *Contact, the contact to send the RPC to
*/
func (network *Network) SendFindDataMessage(hash *KademliaID, contact *Contact) string {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return "ERROR"
	}
	defer conn.Close()
	message := network.getFindDataMessage(hash)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return "ERROR"
	}
	return network.handleSendFindDataResponse(buffer[:n])
}

/*
Returns a standard find data message to find data with hash
  - hash *KademliaID, the hash to find data for
*/
func (network *Network) getFindDataMessage(hash *KademliaID) []byte {
	body, _ := json.Marshal(hash)
	startMessage := []byte(newFindData().startMessage + ";" + string(body) + ";")
	body2 := network.marshalCurrentNode()
	return append(startMessage, body2...)

}

/*
Handle the send find data response and parse to see if an error has occurred or if the response was successful
  - message []byte, the response to parse
*/
func (network *Network) handleSendFindDataResponse(message []byte) string {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return string(message)
	} else {
		if string(message[:4]) == "VALU" {
			resMessage := strings.Split(string(message[5:]), ";")
			var contact *Contact
			json.Unmarshal([]byte(resMessage[1]), &contact)
			if VerifyContact(contact, network) {
				network.RoutingTable.AddContact(*contact)
			}
			return resMessage[0]
		}
		var contacts []Contact
		json.Unmarshal(message, &contacts)
		for _, contact := range contacts {
			if VerifyContact(&contact, network) {
				network.SendPingMessage(&contact)
			}
		}
		return ""
	}
}

/*
Sending a send store RPC message to a contact and returns if the RPC was successful
  - data []byte, the data to store
  - contact *Contact, the contact to send the RPC to
  - kademlia *Kademlia, the kademlia algorithm to send the rpc to
*/
func (network *Network) SendStoreMessage(data []byte, ttl time.Duration, contact *Contact, kademlia *Kademlia) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return false
	}
	defer conn.Close()
	message := network.getStoreMessage(data, ttl)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	hash := NewKademliaID(string(data))
	kademlia.AddToKnown(contact, hash)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	network.handleStoreResponse(buffer[:n])
	return true
}

/*
Returns a standard store message to send the data to store
  - data []byte, the data to send with the store RPC
  - ttl time.Duration, the time to live for the data
*/
func (network *Network) getStoreMessage(data []byte, ttl time.Duration) []byte {
	body, _ := json.Marshal(data)
	body2, _ := json.Marshal(ttl)
	startMessage := []byte(newStoreMessage().startMessage + ";" + string(body) + ";" + string(body2) + ";")
	body3 := network.marshalCurrentNode()
	return append(startMessage, body3...)
}

/*
Handle the store data response and parse to see if an error has occurred or if the response was successful
  - message []byte, the response to parse
*/
func (network *Network) handleStoreResponse(message []byte) {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return
	} else {
		var contact *Contact
		json.Unmarshal(message, &contact)
		if VerifyContact(contact, network) {
			network.RoutingTable.AddContact(*contact)
		}
	}
}

/*
Sending a refresh RPC message to a contact and returns if the RPC was successful
  - hash *KademliaID, the hash of the data to refresh
  - contact *Contact, the contact to send the RPC to
*/
func (network *Network) SendRefreshMessage(hash *KademliaID, contact *Contact) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return false
	}
	defer conn.Close()
	message := network.getRefreshMessage(hash)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	network.handleRefreshResponse(buffer[:n])
	return true
}

/*
Returns a refresh message to send to the contact
  - hash *KademliaId, the hash of the value to refresh
*/
func (network *Network) getRefreshMessage(hash *KademliaID) []byte {
	body, _ := json.Marshal(hash)
	startMessage := []byte(newRefreshMessage().startMessage + ";" + string(body) + ";")
	body2 := network.marshalCurrentNode()
	return append(startMessage, body2...)

}

/*
Handle the refresh response and parse to see if an error has occurred or if the response was successful
  - message []byte, the response to parse
*/
func (network *Network) handleRefreshResponse(message []byte) {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return
	} else {
		var contact *Contact
		json.Unmarshal(message, &contact)
		if VerifyContact(contact, network) {
			network.RoutingTable.AddContact(*contact)
		}
	}
}

/*
A loop that will send a refresh message to every known holder every second
  - kademlia *Kademlia, the current nodes kademlia to get the known holders
*/
func (network *Network) RefreshLoop(kademlia *Kademlia) {
	for {
		for contact, hash := range kademlia.KnownHolders {
			go kademlia.Network.SendRefreshMessage(&hash, &contact)
		}
		time.Sleep(1 * time.Second)
	}
}

/*
Returns the byte array representation of our current nodes contact
information to be sent with every RPC call
*/
func (network *Network) marshalCurrentNode() []byte {
	body, _ := json.Marshal(network.CurrentNode)
	return body
}
