package kademlia

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

/*
Listen function which starts a listening server on input ip and port.
This method will never return any value
  - ip string, the ip address of the server
  - port int, the port number of the server
  - kademlia *Kademlia, the kademlia for this node
*/
func (Network *Network) Listen(ip string, port int, kademlia *Kademlia) {
	addrStr := ip + ":" + strconv.Itoa(port)
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	ln, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println("error is", err)
		return
	}
	fmt.Println("UDP server up and listening on", addrStr)
	defer ln.Close()
	for {
		// wait for UDP client to connect
		handleUDPConnection(ln, kademlia)
	}
}

/*
Handle udp connection by reading bytes and getting the appropriate response message
  - conn *net.UDPConn, the connection to the other node
  - network *Network, the network for this node
  - kademlia *kademlia, the kademlia for this node
*/
func handleUDPConnection(conn *net.UDPConn, kademlia *Kademlia) {
	buffer := make([]byte, maxBytes)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}
	message := getResponseMessage(buffer[:n], kademlia)
	_, err = conn.WriteToUDP(message, addr)
	if err != nil {
		log.Fatal(err)
	}
}

/*
Returns the appropriate response message from a given rpc and calls
the correct functions for each message
  - message []string, the message to decode
  - network *Network, the network our node is on
  - kademlia *Kademlia, the kademlia node we are on
*/
func getResponseMessage(message []byte, kademlia *Kademlia) []byte {
	resMessage := strings.Split(string(message), ";")
	if resMessage[0] == newPing().startMessage {
		return pingResponseMessage(resMessage, kademlia.Network)
	} else if resMessage[0] == newFindContact().startMessage {
		return newFindContactResponseMessage(resMessage, kademlia.Network)
	} else if resMessage[0] == newStoreMessage().startMessage {
		return newStoreMessageResponseMessage(resMessage, kademlia)
	} else if resMessage[0] == newFindData().startMessage {
		return newFindDataResponseMessage(resMessage, kademlia)
	} else if resMessage[0] == newRefreshMessage().startMessage {
		return newRefreshMessageResponseMessage(resMessage, kademlia)
	} else {
		return []byte("Error: Invalid RPC protocol")
	}
}

/*
Returns the response message for when a ping message is received
  - message []string, the message to decode
  - network *Network, the network our node is on
*/
func pingResponseMessage(message []string, network *Network) []byte {
	body := network.marshalCurrentNode()
	ex := extractContact([]byte(message[1]), network)
	if ex != nil {
		return ex
	}
	return body
}

/*
Returns the response message for when a find contact message is received
  - message []string, the message to decode
  - network *Network, the network our node is on
*/
func newFindContactResponseMessage(message []string, network *Network) []byte {
	var id *KademliaID
	json.Unmarshal([]byte(message[1]), &id)
	ex := extractContact([]byte(message[2]), network)
	if ex != nil {
		fmt.Println(ex)
		return ex
	}
	closestNodes := network.RoutingTable.FindClosestContacts(id, BucketSize)
	closestNodes = append(closestNodes, *network.CurrentNode)
	body, _ := json.Marshal(closestNodes)
	return body
}

/*
Returns the response message for when a store message is received
  - message []string, the message to decode
  - network *Network, the network our node is on
  - kademlia *Kademlia, the kademlia node we are on
*/
func newStoreMessageResponseMessage(message []string, kademlia *Kademlia) []byte {
	var data *[]byte
	json.Unmarshal([]byte(message[1]), &data)
	var ttl time.Duration
	json.Unmarshal([]byte(message[2]), &ttl)
	kademlia.Store(*data, ttl)
	ex := extractContact([]byte(message[3]), kademlia.Network)
	if ex != nil {
		fmt.Println(ex)
		return ex
	}
	body := kademlia.Network.marshalCurrentNode()
	return body
}

/*
Returns the response message for when a find data message is received
  - message []string, the message to decode
  - network *Network, the network our node is on
  - kademlia *Kademlia, the kademlia node we are on
*/
func newFindDataResponseMessage(message []string, kademlia *Kademlia) []byte {
	var hash *KademliaID
	json.Unmarshal([]byte(message[1]), &hash)
	ex := extractContact([]byte(message[2]), kademlia.Network)
	if ex != nil {
		fmt.Println(ex)
		return ex
	}
	val := kademlia.LookupData(*hash)
	if val != nil {
		body := kademlia.Network.marshalCurrentNode()
		return []byte("VALU;" + string(val) + ";" + string(body))
	}
	closestNodes := kademlia.Network.RoutingTable.FindClosestContacts(hash, BucketSize)
	closestNodes = append(closestNodes, *kademlia.Network.CurrentNode)
	body, _ := json.Marshal(closestNodes)
	return []byte("CONT" + string(body))
}

/*
Returns the response message for when a refresh message is received
  - message []string, the message to decode
  - network *Network, the network our node is on
  - kademlia *Kademlia, the kademlia node we are on
*/
func newRefreshMessageResponseMessage(message []string, kademlia *Kademlia) []byte {
	var hash *KademliaID
	json.Unmarshal([]byte(message[1]), &hash)
	kademlia.RefreshTTL(*hash)
	ex := extractContact([]byte(message[2]), kademlia.Network)
	if ex != nil {
		fmt.Println(ex)
		return ex
	}
	body := kademlia.Network.marshalCurrentNode()
	return body
}

/*
Verifies if a contact is valid, it is not valid if any is true:
  - is nil
  - is empty string
  - contact id is nil
  - contact id is equal to our id

Returns true if the contact information is correct otherwise false
*/
func VerifyContact(contact *Contact, network *Network) bool {
	return !(contact == nil || contact.Address == "" || contact.ID == nil || contact.ID.Equals(network.CurrentNode.ID))
}

/*
Extracts contact information from a byte array and adds the contact to our routing table,

  - message byte[], the message to extract from
  - network *Network, the network to add the extracted contact to if contact exists

returns an error message if the contact cannot be extracted
*/
func extractContact(message []byte, network *Network) []byte {
	var contact *Contact
	json.Unmarshal(message, &contact)
	if !VerifyContact(contact, network) {
		return []byte("Error: Invalid contact information")
	}
	network.RoutingTable.AddContact(*contact)
	return nil
}
