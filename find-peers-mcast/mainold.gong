package findPeersMCast

import (
	"bufio"
	"fmt"
	"net"
)

type Service struct {
	Denom        string
	Rate         int
	TunnelIP     net.IP
	TunnelPort   int
	TunnelPubkey string
}

func Advertise(
	service *Service,
) error {
	// json, err := json.Marshal(service)
	// if err != nil {
	// 	return err
	// }

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

	fmt.Println("conn created")

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		fmt.Println("scanning")
		text := scanner.Text()
		fmt.Println(text)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func GetPeers() error {
	conn, err := net.DialUDP("udp6", nil, &net.UDPAddr{
		IP:   net.ParseIP("ff02::1"),
		Port: 5544,
		Zone: "eth0",
	})
	if err != nil {
		return err
	}

	conn.Write([]byte("service request\n"))

	// scanner := bufio.NewScanner(conn)

	// for scanner.Scan() {
	// 	text := scanner.Text()
	// 	fmt.Println(text)
	// }
	// if err := scanner.Err(); err != nil {
	// 	return err
	// }

	return nil
}
