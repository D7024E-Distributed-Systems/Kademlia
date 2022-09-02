package network

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

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
		log.Fatal(err)
	}

}
