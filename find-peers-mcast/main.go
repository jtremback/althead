package findNeighborsMCast

import (
	"errors"
	"log"
	"net"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/jtremback/althea/access"
	"github.com/jtremback/althea/serialization"
	"github.com/jtremback/althea/types"
)

type NeighborAPI struct {
	DB *bolt.DB
}

// McastListen listens on the multicast UDP address on a given interface. When it gets
// an althea_hello packet, it calls HelloHandler and sends an althea_hello packet
// to the ControlAddress of the Neighbor.
func (a *NeighborAPI) McastListen(
	iface *net.Interface,
	mcastPort int,
	account types.Account,
	cb func(*types.Neighbor, error),
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
		var b []byte
		_, _, err := conn.ReadFromUDP(b)
		if err != nil {
			cb(nil, err)
			continue
		}

		msg := strings.Split(string(b), " ")

		log.Println("received: " + string(b))

		if msg[0] == "althea_hello" {
			neighbor, err := a.HelloHandler(msg)
			if err != nil {
				cb(nil, err)
				continue
			}

			addr, err := net.ResolveUDPAddr("udp6", neighbor.ControlAddress)
			if err != nil {
				cb(nil, err)
				continue
			}

			err = sendUDP(addr, serialization.FmtHello(account))
			if err != nil {
				cb(nil, err)
				continue
			}

			cb(neighbor, nil)
		} else {
			cb(nil, errors.New("unrecognized message type"))
			continue
		}
	}

	return nil
}

// HelloHandler takes an althea_hello packet, verifies the signature,
// parses the packet into a Neighbor struct, and updates the Neighbor on file with the new information
// contained therein. It also updates the tunneling software. It returns the parsed Neighbor.
// TODO: actually update tunneling software
func (a *NeighborAPI) HelloHandler(msg []string) (*types.Neighbor, error) {
	neighbor, err := serialization.ParseHello(msg)
	if err != nil {
		return nil, err
	}

	err = a.DB.Update(func(tx *bolt.Tx) error {
		return access.SetNeighbor(tx, neighbor)
	})
	if err != nil {
		return nil, err
	}

	return neighbor, nil
}

// ControlListen listens on the ControlAddress of a given interface and passes received messages to
// the appropriate handler function.
func (a *NeighborAPI) ControlListen(
	account types.Account,
	cb func(*types.Neighbor, error),
) error {
	addr, err := net.ResolveUDPAddr("udp6", account.ControlAddress)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp6", addr)
	if err != nil {
		cb(nil, err)
	}

	defer conn.Close()

	for {
		var b []byte
		_, _, err := conn.ReadFromUDP(b)
		if err != nil {
			cb(nil, err)
			continue
		}

		msg := strings.Split(string(b), " ")

		log.Println("received: " + string(b))

		if msg[0] == "althea_hello" {
			_, err := a.HelloHandler(msg)
			if err != nil {
				cb(nil, err)
				continue
			}
		} else {
			cb(nil, errors.New("unrecognized message type"))
			continue
		}
	}

	return nil
}

// McastHello sends an althea_hello packet to the multicast UDP address on a given interface.
func McastHello(
	iface *net.Interface,
	mCastPort int,
	account types.Account,
	cb func(*types.Neighbor, error),
) error {
	err := sendUDP(&net.UDPAddr{
		IP:   net.ParseIP("ff02::1"),
		Port: mCastPort,
		Zone: iface.Name,
	}, serialization.FmtHello(account))
	if err != nil {
		return err
	}
	return nil
}

func sendUDP(
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
