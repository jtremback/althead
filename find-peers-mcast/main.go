package findPeersMCast

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type Service struct {
	Denom        string
	Rate         int
	TunnelIP     net.IP
	TunnelPort   int
	TunnelPubkey string
}

func Advertise() error {
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		return err
	}

	conn, err := net.ListenMulticastUDP(
		"udp6",
		iface,
		&net.UDPAddr{
			IP:   net.ParseIP("ff02::1"),
			Port: 5544,
			Zone: "eth0",
		},
	)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(conn)
	fmt.Println("made conn")
	for scanner.Scan() {
		fmt.Println("scannin")
		text := scanner.Text()
		fmt.Println(text)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func QueryPeers(
	IP net.IP,
	port int,
) error {
	conn, err := net.DialUDP(
		"udp6",
		nil,
		&net.UDPAddr{
			IP:   net.ParseIP("ff02::1"),
			Port: 5544,
			Zone: "eth0",
		})
	if err != nil {
		return err
	}

	conn.Write([]byte("Althea service request " +
		IP.String() +
		" " +
		strconv.Itoa(port) +
		"\n"))

	return nil
}
