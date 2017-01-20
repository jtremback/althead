package neighborAPI

import (
	"errors"
	"log"
	"net"
	"strings"

	"github.com/agl/ed25519"
	"github.com/jtremback/scrooge/serialization"
	"github.com/jtremback/scrooge/types"
)

type NeighborAPI struct {
	Neighbors map[[ed25519.PublicKeySize]byte]*types.Neighbor
	Account   *types.Account
}

// McastListen listens on the multicast UDP address on a given interface. When it gets
// an scrooge_hello packet, it calls HelloHandler and sends an scrooge_hello packet
// to the ControlAddress of the Neighbor.
func (a *NeighborAPI) McastListen(
	mcastPort int,
	iface *net.Interface,
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

		if msg[0] == "scrooge_hello" {
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

			err = sendUDP(addr, serialization.FmtHello(a.Account))
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

func (a *NeighborAPI) createTunnel() {

}

func (a *NeighborAPI) HelloHandler(msg []string) (*types.Neighbor, error) {
	helloMessage, err := serialization.ParseHello(msg)
	if err != nil {
		return nil, err
	}

	neighbor := a.Neighbors[helloMessage.PublicKey]
	if neighbor == nil {
		neighbor = &types.Neighbor{
			PublicKey: helloMessage.PublicKey,
		}
	}

	if neighbor.Seqnum > helloMessage.Seqnum {
		return nil, errors.New("sequence number too low")
	}

	neighbor.Seqnum = helloMessage.Seqnum
	neighbor.ControlAddress = helloMessage.ControlAddress

	return neighbor, nil
}

// ControlListen listens on the ControlAddress and passes received messages to
// the appropriate handler function.
func (a *NeighborAPI) ControlListen(
	cb func(*types.Neighbor, error),
) error {
	addr, err := net.ResolveUDPAddr("udp6", a.Account.ControlAddress)
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

		if msg[0] == "scrooge_hello" {
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

// McastHello sends an scrooge_hello packet to the multicast UDP address on a given interface.
func (a *NeighborAPI) McastHello(
	mCastPort int,
	iface *net.Interface,
	cb func(*types.Neighbor, error),
) error {
	err := sendUDP(&net.UDPAddr{
		IP:   net.ParseIP("ff02::1"),
		Port: mCastPort,
		Zone: iface.Name,
	}, serialization.FmtHello(a.Account))
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
