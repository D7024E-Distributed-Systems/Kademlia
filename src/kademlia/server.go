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
		handleUDPConnection(ln, Network, kademlia)
	}
}

func handleUDPConnection(conn *net.UDPConn, Network *Network, kademlia *Kademlia) {

	// here is where you want to do stuff like read or write to client

	buffer := make([]byte, maxBytes)
	n, addr, err := conn.ReadFromUDP(buffer)

	if err != nil {
		log.Fatal(err)
	}

	message := getResponseMessage(buffer[:n], Network, kademlia)

	_, err = conn.WriteToUDP(message, addr)

	if err != nil {
		log.Fatal(err)
	}
}

func getResponseMessage(message []byte, Network *Network, kademlia *Kademlia) []byte {
	resMessage := strings.Split(string(message), ";")
	if resMessage[0] == newPing().startMessage {
		body := Network.marshalCurrentNode()
		ex := extractContact([]byte(resMessage[1]), Network)
		if ex != nil {
			return ex
		}
		return body
	} else if resMessage[0] == newFindContact().startMessage {
		var id *KademliaID
		json.Unmarshal([]byte(resMessage[1]), &id)
		ex := extractContact([]byte(resMessage[2]), Network)
		if ex != nil {
			fmt.Println(ex)
			return ex
		}
		closestNodes := Network.RoutingTable.FindClosestContacts(id, BucketSize)
		closestNodes = append(closestNodes, *Network.CurrentNode)
		body, _ := json.Marshal(closestNodes)
		return body
	} else if resMessage[0] == newStoreMessage().startMessage {
		var data *[]byte
		json.Unmarshal([]byte(resMessage[1]), &data)
		var ttl time.Duration
		json.Unmarshal([]byte(resMessage[2]), &ttl)
		kademlia.Store(*data, ttl)
		ex := extractContact([]byte(resMessage[3]), Network)
		if ex != nil {
			fmt.Println(ex)
			return ex
		}
		body := Network.marshalCurrentNode()
		return body
	} else if resMessage[0] == newFindData().startMessage {
		var hash *KademliaID
		json.Unmarshal([]byte(resMessage[1]), &hash)
		ex := extractContact([]byte(resMessage[2]), Network)
		if ex != nil {
			fmt.Println(ex)
			return ex
		}
		val := kademlia.LookupData(*hash)
		if val != nil {
			body := Network.marshalCurrentNode()
			return []byte("VALU;" + string(val) + ";" + string(body))
		}
		closestNodes := Network.RoutingTable.FindClosestContacts(hash, BucketSize)
		closestNodes = append(closestNodes, *Network.CurrentNode)
		body, _ := json.Marshal(closestNodes)
		return []byte("CONT" + string(body))
	} else if resMessage[0] == newRefreshmessage().startMessage {
		var hash *KademliaID
		json.Unmarshal([]byte(resMessage[1]), &hash)
		kademlia.RefreshTTL(*hash)
		ex := extractContact([]byte(resMessage[2]), Network)
		if ex != nil {
			fmt.Println(ex)
			return ex
		}
		body := Network.marshalCurrentNode()
		return body
	} else {
		return []byte("Error: Invalid RPC protocol")
	}

}

/**
 * returns true if the contact information is correct
 */
func VerifyContact(contact *Contact, network *Network) bool {
	return !(contact == nil || contact.Address == "" || contact.ID == nil || contact.ID.Equals(network.CurrentNode.ID))
}

func extractContact(message []byte, network *Network) []byte {
	var contact *Contact
	json.Unmarshal(message, &contact)
	if !VerifyContact(contact, network) {
		return []byte("Error: Invalid contact information")
	}
	network.RoutingTable.AddContact(*contact)
	return nil
}
