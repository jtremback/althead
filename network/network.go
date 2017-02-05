package network

import (
	"log"
	"net"
)

type Network struct{}

// McastListen listens on the multicast UDP address on a given interface.
func (self *Network) McastListen(
	port int,
	iface *net.Interface,
	handlers func([]byte, *net.Interface) error,
	cb func(error),
) error {
	conn, err := net.ListenMulticastUDP(
		"udp6",
		iface,
		&net.UDPAddr{
			IP:   net.ParseIP("ff02::1"),
			Port: port,
			Zone: iface.Name,
		},
	)
	if err != nil {
		return err
	}

	for {
		var b []byte
		_, _, err := conn.ReadFromUDP(b)
		if err != nil {
			cb(err)
			continue
		}
		cb(handlers(b, iface))
	}

	return nil
}

func (self *Network) UnicastListen(
	addr *net.UDPAddr,
	iface *net.Interface,
	handlers func([]byte, *net.Interface) error,
	cb func(error),
) error {
	conn, err := net.ListenUDP("udp6", addr)
	if err != nil {
		return err
	}

	defer conn.Close()

	for {
		var b []byte
		_, _, err := conn.ReadFromUDP(b)
		if err != nil {
			cb(err)
			continue
		}
		cb(handlers(b, iface))
	}

	return nil
}

func (self *Network) SendUDP(
	addr *net.UDPAddr,
	s string,
) error {
	conn, err := net.DialUDP(
		"udp6",
		nil,
		addr,
	)
	defer conn.Close()
	if err != nil {
		return err
	}

	conn.Write([]byte(s))
	log.Println("sent: " + s)
	return nil
}

// McastHello sends a packet to the multicast UDP address on a given interface.
func (self *Network) SendMulticastUDP(
	iface *net.Interface,
	port int,
	s string,
) error {
	err := self.SendUDP(&net.UDPAddr{
		IP:   net.ParseIP("ff02::1"),
		Port: port,
		Zone: iface.Name,
	}, s)
	if err != nil {
		return err
	}
	return nil
}
