package d7024e

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

type Network struct {
	currentNode *Contact
}

func NewNetwork(node *Contact) *Network {
	return &Network{node}
}

func Listen(ip string, port int) {
	addrStr := ip + ":" + strconv.Itoa(port)
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(addr)
	ln, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println("error is", err)
		return
	}

	fmt.Println("UDP server up and listening on port", port)

	defer ln.Close()

	for {
		// wait for UDP client to connect
		handleUDPConnection(ln)
	}
}

func handleUDPConnection(conn *net.UDPConn) {

	// here is where you want to do stuff like read or write to client

	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP client :", addr)
	fmt.Println("Received from UDP client :", string(buffer[:n]))

	// TODO: Call correct method depending on received message and reply

	// write message back to client
	message := []byte("Hello UDP client!")
	_, err = conn.WriteToUDP(message, addr)

	if err != nil {
		log.Println(err)
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Println("SendPingMessage")
	fmt.Println(contact.Address)

	conn, err3 := net.Dial("udp4", contact.Address)
	defer conn.Close()

	if err3 != nil {
		log.Println(err3)
	}
	message := []byte("Hello UDP server!")
	conn.Write(message)
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(buffer[:n]))

	// TODO
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
