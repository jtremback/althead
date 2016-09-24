package findPeersMCast

import (
	"errors"
	"fmt"
	"net"
	"time"
)

func Advertise(
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
		b := make([]byte, 12)
		_, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			return err
		}
		if string(b) == "althea_hello" {
			conn, err := net.DialUDP(
				"udp6",
				nil,
				addr,
			)
			if err != nil {
				return err
			}

			conn.Write([]byte("althea_ihu"))

			conn.Close()

		}
		fmt.Println("got msg", string(b), addr)
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

func QueryPeers(
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

	go func() {
		for {
			b := make([]byte, 10)
			_, addr, err := conn.ReadFromUDP(b)
			if err != nil {
				cb(nil, err)
			}
			fmt.Println(addr.String())
			cb(net.ParseIP(addr.String()), nil)
		}
	}()

	conn.WriteToUDP([]byte("althea_hello"), &net.UDPAddr{
		IP:   net.ParseIP("ff02::1"),
		Port: mCastPort,
		Zone: iface.Name,
	})

	time.Sleep(1 * time.Second)
}
