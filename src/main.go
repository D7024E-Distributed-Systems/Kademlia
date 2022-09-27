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

/*
 * Starts the program by connecting to the kademlia network and start listening services
 */
func main() {
	// Default ip and port for first connection to Kademlia network
	port := 3000
	restPort := 3001
	// defaultIp := "192.168.1.182"
	defaultIp := "172.19.0.2"
	// defaultIp := "173.19.0.2"

	// sets a random seed for the random number generator
	rand.Seed(time.Now().UnixNano())

	// sets up a contact to be able to send ping and later be changed
	currentContact := NewContact(NewRandomKademliaID(), "")
	// get the ip of the current node
	ip := getOutboundIP()
	var kademlia *Kademlia
	var contact Contact
	// if the defaultIp is equal to the nodes ip we know that we will start this node
	// wether or not we can connect to the network
	if ip.String() == defaultIp {
		// sends a ping to the network defaultIp and defaultPort
		contact = NewContact(NewRandomKademliaID(), defaultIp+":"+strconv.Itoa(port))
		success := NewNetwork(&currentContact).SendPingMessage(&contact)
		// if the response was successful then we set a random port for this node
		if success {
			port = rand.Intn(65535-1024) + 1024
			restPort = rand.Intn(65535-1024) + 1024
		}
		// sets the currentContact and the kademlia network
		currentContact, kademlia = createCurrentContact(ip, port)
	} else {
		// if we're at a different ip address then the defaultIp then we try to connect to the network
		contact = NewContact(NewRandomKademliaID(), defaultIp+":"+strconv.Itoa(port))
		currentContact, kademlia = createCurrentContact(ip, port)
		success := kademlia.Network.SendPingMessage(&contact)
		// if not successful we will panic since we couldn't connect to the kademlia network
		if !success {
			fmt.Println("Our IP address is", ip.String())
			panic("failed to connect to p2p server")
		}
		// sets a random port number
		port = rand.Intn(65535-1024) + 1024
		restPort = rand.Intn(65535-1024) + 1024
	}
	fmt.Println("Current contact main", currentContact)
	// start this node to listen to incoming messages
	go kademlia.Network.Listen(ip.String(), port, kademlia)
	go GetRoute(ip.String(), restPort, kademlia)
	// starts the cli "class"
	go cli.Init(shutdownNode, kademlia)
	// sleep for 1 second to make sure everything has time to be initialized
	time.Sleep(1 * time.Second)
	// call findContactMessage for find ourself
	kademlia.Network.SendFindContactMessage(&contact, currentContact.ID)
	i := 0
	// for loop to not exit the main thread
	for {
		i++
		kademlia.DeleteOldData()
		for contact, hash := range kademlia.KnownHolders {
			go kademlia.Network.SendRefreshMessage(&hash, &contact)
		}
		if i%6 == 0 {
			fmt.Println(kademlia.Network.RoutingTable.FindClosestContacts(NewRandomKademliaID(), 1000))
		}
		time.Sleep(5 * time.Second)
	}
}

/*
 * returns new contact and kademlia instance for the given ip and port
 */
func createCurrentContact(ip net.IP, port int) (Contact, *Kademlia) {
	contact := NewContact(NewRandomKademliaID(), ip.String()+":"+strconv.Itoa(port))
	kademlia := NewKademliaStruct(NewNetwork(&contact))
	return contact, kademlia
}

/*
 * Returns the ip of the current node
 */
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

/*
 * Shuts down the node
 */
func shutdownNode() {
	fmt.Println("Shutting down node")
	os.Exit(0)
}
