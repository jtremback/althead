package findPeersMCast

import (
	"errors"
	"log"
	"net"
	"strings"
)

func Listen(
	iface *net.Interface,
	mcastPort int,
) error {

	conn, err := net.ListenMulticastUDP(
		"udp6",
		iface,
		&net.UDPAddr{
			IP:   net.ParseIP("ff02::1"),
			Port: mcastPort,
			Zone: iface.Name,
		},
	)
	if err != nil {
		return err
	}

	for {
		b := make([]byte, 64)
		_, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			return err
		}

		msg := strings.Split(string(b), " ")

		log.Println("received: "+string(b), "from: ", addr)

		if msg[0] == "althea_hello" {
			conn, err := net.DialUDP(
				"udp6",
				nil,
				addr,
			)
			if err != nil {
				return err
			}

			s := "althea_ihu <pubkey>"

			conn.Write([]byte(s))

			log.Println("sent: " + s)
			conn.Close()

		}
	}

	return nil
}

func firstLinkLocalUnicast(iface *net.Interface) (*net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		}

		if ip.IsLinkLocalUnicast() {
			return &ip, nil
		}
	}
	return nil, errors.New("Could not find link local unicast ipv6 address for interface " + iface.Name)
}

func Hello(
	iface *net.Interface,
	mCastPort int,
	cb func(net.IP, error),
) {
	ip, err := firstLinkLocalUnicast(iface)
	if err != nil {
		cb(nil, err)
	}

	conn, err := net.ListenUDP("udp6", &net.UDPAddr{
		IP:   *ip,
		Port: 0,
		Zone: iface.Name,
	})
	if err != nil {
		cb(nil, err)
	}

	defer conn.Close()

	s := "althea_hello <pubkey>"

	conn.WriteToUDP([]byte(s), &net.UDPAddr{
		IP:   net.ParseIP("ff02::1"),
		Port: mCastPort,
		Zone: iface.Name,
	})

	log.Println("sent: " + s)

	for {
		b := make([]byte, 64)
		_, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			cb(nil, err)
		}
		log.Println("received: "+string(b), "from: ", addr)
		cb(net.ParseIP(addr.String()), nil)
	}
}
