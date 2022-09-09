package network

import (
	"encoding/json"
	"log"
	"net"
	"time"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
)

/**
 * ping = PING
 * find contact =FICO
 * find data = FIDA
 * store message = STME
 */
var maxBytes int = 4096

func (network *Network) SendPingMessage(contact *Contact) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
	}
	defer conn.Close()
	message := getPingMessage(network)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	// fmt.Println("\tResponse from server:", string(buffer[:n]))
	handlePingResponse(buffer[:n], network)
	return true
}

func getPingMessage(network *Network) []byte {
	startMessage := []byte(newPing().startMessage + ";")
	body := network.marshalCurrentNode()
	return append(startMessage, body...)
}

func handlePingResponse(message []byte, network *Network) {
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
	// fmt.Println("ping response: ", network.routingTable)
}

func (network *Network) SendFindContactMessage(contact *Contact, nodeID *KademliaID) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
	}
	defer conn.Close()
	message := getFindContactMessage(network, nodeID)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	// fmt.Println("\tResponse from server:", string(buffer[:n]))
	handleFindContactResponse(buffer[:n], network)
	return true
	// TODO
}

func getFindContactMessage(network *Network, nodeID *KademliaID) []byte {
	body, err := json.Marshal(nodeID)
	if err != nil {
		log.Println(err)
	}
	startMessage := []byte(newFindContact().startMessage + ";" + string(body) + ";")
	body2 := network.marshalCurrentNode()
	return append(startMessage, body2...)

}

func handleFindContactResponse(message []byte, network *Network) {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return
	} else {
		var contacts []Contact
		json.Unmarshal(message, &contacts)
		for _, contact := range contacts {
			if VerifyContact(&contact, network) {
				network.SendPingMessage(&contact)
			}
		}
	}
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte, contact *Contact) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
	}
	defer conn.Close()
	message := getStoreMessage(network, data)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	// fmt.Println("\tResponse from server:", string(buffer[:n]))
	handleStoreResponse(buffer[:n], network)
	return true
	// TODO
}

func getStoreMessage(network *Network, data []byte) []byte {
	body, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	startMessage := []byte(newStoreMessage().startMessage + ";" + string(body) + ";")
	body2 := network.marshalCurrentNode()
	return append(startMessage, body2...)
}

func handleStoreResponse(message []byte, network *Network) {
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

func (network *Network) marshalCurrentNode() []byte {
	body, err := json.Marshal(network.CurrentNode)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return body
}
