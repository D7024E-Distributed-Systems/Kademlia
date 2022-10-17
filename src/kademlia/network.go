package kademlia

import (
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"
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
		return false
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
}

func (network *Network) SendFindContactMessage(contact *Contact, nodeID *KademliaID) []Contact {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return nil
	}
	defer conn.Close()
	message := getFindContactMessage(network, nodeID)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return nil
	}
	return handleFindContactResponse(buffer[:n], network)
}

func getFindContactMessage(network *Network, nodeID *KademliaID) []byte {
	body, _ := json.Marshal(nodeID)
	startMessage := []byte(newFindContact().startMessage + ";" + string(body) + ";")
	body2 := network.marshalCurrentNode()
	return append(startMessage, body2...)

}

func handleFindContactResponse(message []byte, network *Network) []Contact {
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

func (network *Network) SendFindDataMessage(hash *KademliaID, contact *Contact) string {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return "ERROR"
	}
	defer conn.Close()
	message := getFindDataMessage(network, hash)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return "ERROR"
	}
	return handleSendDataResponse(buffer[:n], network)

	// TODO
}

func getFindDataMessage(network *Network, hash *KademliaID) []byte {
	body, _ := json.Marshal(hash)
	startMessage := []byte(newFindData().startMessage + ";" + string(body) + ";")
	body2 := network.marshalCurrentNode()
	return append(startMessage, body2...)

}

func handleSendDataResponse(message []byte, network *Network) string {
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

func (network *Network) SendStoreMessage(data []byte, ttl time.Duration, contact *Contact, kademlia *Kademlia) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return false
	}
	defer conn.Close()
	message := getStoreMessage(network, data, ttl)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	hash := NewKademliaID(string(data))
	kademlia.AddToKnown(contact, hash)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	handleStoreResponse(buffer[:n], network)
	return true
}

func getStoreMessage(network *Network, data []byte, ttl time.Duration) []byte {
	body, _ := json.Marshal(data)
	body2, _ := json.Marshal(ttl)
	startMessage := []byte(newStoreMessage().startMessage + ";" + string(body) + ";" + string(body2) + ";")
	body3 := network.marshalCurrentNode()
	return append(startMessage, body3...)
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

func (network *Network) SendRefreshMessage(hash *KademliaID, contact *Contact) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return false
	}
	defer conn.Close()
	message := getRefreshMessage(network, hash)
	conn.Write(message)
	buffer := make([]byte, maxBytes)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	handleRefreshResponse(buffer[:n], network)
	return true
}

func getRefreshMessage(network *Network, hash *KademliaID) []byte {
	body, _ := json.Marshal(hash)
	startMessage := []byte(newRefreshMessage().startMessage + ";" + string(body) + ";")
	body2 := network.marshalCurrentNode()
	return append(startMessage, body2...)

}

func (network *Network) RefreshLoop(kademlia *Kademlia) {
	for {
		for contact, hash := range kademlia.KnownHolders {
			go kademlia.Network.SendRefreshMessage(&hash, &contact)
		}
		time.Sleep(1 * time.Second)
	}
}

func handleRefreshResponse(message []byte, network *Network) {
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
	body, _ := json.Marshal(network.CurrentNode)
	return body
}
