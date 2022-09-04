package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/d7024e"
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/network"
)

func main() {
	// Default ip and port for first connection to Kademlia network
	port := 3000
	defaultIp := "130.240.108.230"
	// The target is just some random ID and default ip and port
	target := NewRandomKademliaID()
	contact := NewContact(target, defaultIp+":"+strconv.Itoa(port))
	// Our current contact, which is this node, will be some random ID and no address
	currentContact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&currentContact)
	// we store the success of the ping message, if the ping was successful then we can
	// start our server on a random port number on our local network. Else we start ourself
	// at default port.
	success := network.SendPingMessage(&contact)
	ip := getOutboundIP()
	if success {
		rand.Seed(time.Now().UnixNano())
		// random port number
		port = rand.Intn(65535-1024) + 1024
	}
	currentContact, network = createCurrentContact(ip, port)
	go network.Listen(ip.String(), port)
	// go network.SendFindContactMessage(&currentContact)
	go network.SendPingMessage(&contact)
	fmt.Println("Current contact main", currentContact)
	for {
		fmt.Println(network.RoutingTable.FindClosestContacts(currentContact.ID, 1000))
		time.Sleep(15 * time.Second)
	}
}

func createCurrentContact(ip net.IP, port int) (Contact, *Network) {
	contact := NewContact(NewRandomKademliaID(), ip.String()+":"+strconv.Itoa(port))
	network := NewNetwork(&contact)
	return contact, network
}

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
