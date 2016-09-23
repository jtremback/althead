package findPeersMCast

import (
	"errors"
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
		if string(b) == "althea_hello" {
			conn, err := net.DialUDP(
				"udp6",
				nil,
				addr,
			)
			if err != nil {
				return err
			}

			conn.Write([]byte("althea_hello"))

			conn.Close()

		}
		fmt.Println("foobler", string(b), addr)
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
	listenPort int,
	mCastPort int,
) error {
	ip, err := firstLinkLocalUnicast(iface)
	if err != nil {
		return err
	}

	laddr := &net.UDPAddr{
		IP:   *ip,
		Port: 0,
		Zone: iface.Name,
	}
	fmt.Println(laddr)
	l, err := net.ListenUDP("udp6", laddr)
	if err != nil {
		return err
	}

	laddr, err = net.ResolveUDPAddr("udp6", l.LocalAddr().String())
	if err != nil {
		return err
	}
	fmt.Println(laddr)
	defer l.Close()

	ch := make(chan error)

	go func() {
		for {
			_, addr, err := l.ReadFromUDP([]byte{})
			if err != nil {
				ch <- err
			}
			fmt.Println("goobler", addr)
		}
	}()

	raddr := &net.UDPAddr{
		IP:   net.ParseIP("ff02::1"),
		Port: mCastPort,
		Zone: iface.Name,
	}

	conn, err := net.DialUDP("udp6", nil, raddr)
	if err != nil {
		return err
	}

	conn.Write([]byte("althea_hello " +
		laddr.String()))

	conn.Close()

	err = <-ch
	return err
}
