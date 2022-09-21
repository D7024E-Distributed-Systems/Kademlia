package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/D7024E-Distributed-Systems/Kademlia/src/cli"
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/rest"
)

func main() {
	// Default ip and port for first connection to Kademlia network

	port := 3000
	restPort := 3001
	// defaultIp := "130.240.153.179"
	defaultIp := "172.19.0.2"
	// cli.Init(shutdownNode)

	/** //! UNCOMMENT THIS WHEN WE WANT TO GO TO PRODUCTION
	ip := getOutboundIP()
	var currentContact Contact
	var network *Network
	if ip.String() == defaultIp {
		// If we are at the local network address we just open the port at 3000
		currentContact, network = createCurrentContact(ip, port)
		go network.Listen(ip.String(), port)
	} else {
		target := NewRandomKademliaID()
		contact := NewContact(target, defaultIp+":"+strconv.Itoa(port))
		rand.Seed(time.Now().UnixNano())
		// random port number
		port = rand.Intn(65535-1024) + 1024
		currentContact, network = createCurrentContact(ip, port)
		go network.Listen(ip.String(), port)
		success := network.SendPingMessage(&contact)
		if !success {
			panic("failed to connect to p2p server")
		}
	}
	fmt.Println("Current contact main", currentContact)
	*/

	// The target is just some random ID and default ip and port
	target := NewRandomKademliaID()
	gatekeeper := NewContact(target, defaultIp+":"+strconv.Itoa(port))
	// Our current contact, which is this node, will be some random ID and no address
	currentContact := NewContact(NewRandomKademliaID(), "")
	// we store the success of the ping message, if the ping was successful then we can
	// start our server on a random port number on our local network. Else we start ourself
	// at default port.
	ip := getOutboundIP()
	fmt.Println(ip)
	success := NewNetwork(&currentContact).SendPingMessage(&gatekeeper)
	if ip.String() != defaultIp && !success {
		fmt.Println("Our IP address is", ip.String())
		panic("Couldn't connect to p2p network")
	}
	if success {
		fmt.Println("Random port")
		rand.Seed(time.Now().UnixNano())
		// random port number
		port = rand.Intn(65535-1024) + 1024
		rand.Seed(time.Now().UnixNano())

		restPort = rand.Intn(65535-1024) + 1024
	}
	currentContact, kademlia := createCurrentContact(ip, port)
	if success {
		kademlia.Network.SendStoreMessage([]byte("String"), 15*time.Minute, &gatekeeper, kademlia)
	}
	go kademlia.Network.Listen(ip.String(), port, kademlia)
	// go network.SendFindContactMessage(&currentContact)
	// go network.SendPingMessage(&contact)
	kademlia.Network.SendFindContactMessage(&gatekeeper, currentContact.ID)
	// hash := NewKademliaID("String")
	fmt.Println("Current contact main", currentContact)
	go GetRoute(ip.String(), restPort, kademlia)
	go cli.Init(shutdownNode, kademlia)
	// i := 0
	for {
		// fmt.Println(kademlia.Network.RoutingTable.FindClosestContacts(currentContact.ID, 1000))
		kademlia.DeleteOldData()
		for contact, hash := range kademlia.KnownHolders {
			go kademlia.Network.SendRefreshMessage(&hash, &contact)
		}
		// kademlia.Network.SendStoreMessage([]byte("String"), 15*time.Minute, &gatekeeper, kademlia)
		time.Sleep(5 * time.Second)
		// res := network.SendFindDataMessage(hash, &gatekeeper)
		// fmt.Println("THE SAVED VALUE IS:", res)
	}
}

func createCurrentContact(ip net.IP, port int) (Contact, *Kademlia) {
	contact := NewContact(NewRandomKademliaID(), ip.String()+":"+strconv.Itoa(port))
	kademlia := NewKademliaStruct(NewNetwork(&contact))
	return contact, kademlia
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

func shutdownNode() {
	fmt.Println("Shutting down node")
	os.Exit(0)
}
