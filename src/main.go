package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/d7024e"
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/network"
)

func main() {
	target := NewRandomKademliaID()
	contact := NewContact(target, "127.0.0.1:3000")
	contact.CalcDistance(target)
	currentContact := NewContact(NewRandomKademliaID(), "")
	network := NewNetwork(&currentContact)
	ping := network.SendPingMessage(&contact)
	if ping {
		ip := GetOutboundIP()
		// port := rand.Int(65535-1024) + 1024
		rand.Seed(time.Now().UnixNano())
		port := rand.Intn(65535-1024) + 1024
		fmt.Println(ip)
		fmt.Println(port)
		go Listen(ip.String(), port)
		// TODO: Send find node to server and start listening on random port
	} else {
		go Listen("127.0.0.1", 3000)
	}
	for {
		time.Sleep(2 * time.Second)
	}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
