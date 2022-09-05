package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/d7024e"
)

func (Network *Network) Listen(ip string, port int) {
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
		handleUDPConnection(ln, Network)
	}
}

func handleUDPConnection(conn *net.UDPConn, Network *Network) {

	// here is where you want to do stuff like read or write to client

	buffer := make([]byte, maxBytes)
	n, addr, err := conn.ReadFromUDP(buffer)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\tReceived from UDP client :", string(buffer[:n]))

	message := getResponseMessage(buffer[:n], Network)

	// TODO: Call correct method depending on received message and reply

	// write message back to client
	// message := []byte("Hello UDP client!")
	_, err = conn.WriteToUDP(message, addr)

	if err != nil {
		log.Fatal(err)
	}
}

func getResponseMessage(message []byte, Network *Network) []byte {
	messageCode := string(message[:4])
	if messageCode == newPing().startMessage {
		body, err := json.Marshal(Network.CurrentNode)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		ex := extractContact(message[4:], Network)
		if ex != nil {
			return ex
		}
		return body
	} else if messageCode == newFindContact().startMessage {
		res := strings.Split(string(message[4:]), ";")
		var id *KademliaID
		json.Unmarshal([]byte(res[0]), &id)
		ex := extractContact([]byte(res[1]), Network)
		if ex != nil {
			fmt.Println(ex)
			// return ex
		}
		closestNodes := Network.RoutingTable.FindClosestContacts(id, BucketSize)
		closestNodes = append(closestNodes, *Network.CurrentNode)
		body, err := json.Marshal(closestNodes)
		if err != nil {
			log.Println(err)
			panic(err)
		}
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
