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
	defaultIp := "192.168.1.182"
	// defaultIp := "172.19.0.2"
	rand.Seed(time.Now().UnixNano())
	// defaultIp := "173.19.0.2"

	currentContact := NewContact(NewRandomKademliaID(), "")
	ip := getOutboundIP()
	var kademlia *Kademlia
	var contact Contact
	if ip.String() == defaultIp {
		contact = NewContact(NewRandomKademliaID(), defaultIp+":"+strconv.Itoa(port))
		success := NewNetwork(&currentContact).SendPingMessage(&contact)
		if success {
			port = rand.Intn(65535-1024) + 1024
			restPort = rand.Intn(65535-1024) + 1024
		}
		// If we are at the local network address we just open the port at 3000
		currentContact, kademlia = createCurrentContact(ip, port)
	} else {
		contact = NewContact(NewRandomKademliaID(), defaultIp+":"+strconv.Itoa(port))
		// random port number
		port = rand.Intn(65535-1024) + 1024
		restPort = rand.Intn(65535-1024) + 1024
		currentContact, kademlia = createCurrentContact(ip, port)
		success := kademlia.Network.SendPingMessage(&contact)
		if !success {
			fmt.Println("Our IP address is", ip.String())
			panic("failed to connect to p2p server")
		}
	}
	fmt.Println("Current contact main", currentContact)
	go kademlia.Network.Listen(ip.String(), port, kademlia)
	go GetRoute(ip.String(), restPort, kademlia)
	go cli.Init(shutdownNode, kademlia)
	kademlia.Network.SendFindContactMessage(&contact, currentContact.ID)
	i := 0
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
