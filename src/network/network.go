package network

import (
	"fmt"
	"log"
	"net"
	"time"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/d7024e"
)

func (network *Network) SendPingMessage(contact *Contact) bool {
	fmt.Println(contact.Address)
	conn, err3 := net.Dial("udp4", contact.Address)
	defer conn.Close()
	if err3 != nil {
		log.Println(err3)
	}
	message := []byte("Hello UDP server!")
	conn.Write(message)
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	fmt.Println("Response from server:", string(buffer[:n]))
	return true
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
