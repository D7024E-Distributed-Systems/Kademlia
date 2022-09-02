package d7024e

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

type Network struct {
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

	// var testPayLoad []byte = []byte("This is a test")

	// conn.Write(testPayLoad)
	// TODO
}

func handleUDPConnection(conn *net.UDPConn) {

	// here is where you want to do stuff like read or write to client

	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP client : ", addr)
	fmt.Println("Received from UDP client :  ", string(buffer[:n]))

	if err != nil {
		log.Fatal(err)
	}

	// NOTE : Need to specify client address in WriteToUDP() function
	//        otherwise, you will get this error message
	//        write udp : write: destination address required if you use Write() function instead of WriteToUDP()

	// write message back to client
	message := []byte("Hello UDP client!")
	_, err = conn.WriteToUDP(message, addr)

	if err != nil {
		log.Println(err)
	}

}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {

	fmt.Println(buf)
	buf[2] |= 0x80
	fmt.Println(buf)

	pc.WriteTo(buf, addr)
}

func (network *Network) SendPingMessage(contact *Contact) {

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
