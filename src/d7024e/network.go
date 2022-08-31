package d7024e

import (
	"fmt"
	"net"
)

type Network struct {
}

func Listen(ip string, port int) {
	byteIP := []byte(ip)
	laddr := net.UDPAddr{IP: byteIP, Port: port, Zone: ""}
	conn, err := net.ListenUDP("udp", &laddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var buf [1024]byte

	for {
		rlen, remote, err := conn.ReadFromUDP(buf[:])

		if err != nil {
			continue
		}

		go serve(conn, remote, buf[:rlen])
	}

	var testPayLoad []byte = []byte("This is a test")

	conn.Write(testPayLoad)
	// TODO
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
