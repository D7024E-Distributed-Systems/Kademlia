package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

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

	// fmt.Println("\tReceived from UDP client :", string(buffer[:n]))
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
	if string(message[:4]) == newPing().startMessage {
		var contact *Contact
		json.Unmarshal(message[4:], &contact)
		if !VerifyContact(contact, Network) {
			return []byte("Error: Invalid contact information")
		}
		Network.RoutingTable.AddContact(*contact)
		body, err := json.Marshal(Network.CurrentNode)
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

	return !(contact.Address == "" || contact.ID == nil || contact.ID.Equals(network.CurrentNode.ID))
}
