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
	Network   interface {
		SendUDP(*net.UDPAddr, string) error
	}
}

// McastListen listens on the multicast UDP address on a given interface. When it gets
// an scrooge_hello packet, it calls HelloHandler and sends an scrooge_hello packet
// to the ControlAddress of the Neighbor.
func (self *NeighborAPI) McastListen(
	port int,
	iface *net.Interface,
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
		cb(self.Handlers(b))
	}

	return nil
}

// ControlListen listens on the ControlAddress and passes received messages to
// the appropriate handler function.
func (self *NeighborAPI) ControlListen(
	port int,
	iface *net.Interface,
	cb func(error),
) error {
	addr, err := net.ResolveUDPAddr("udp6", self.Account.ControlAddress)
	if err != nil {
		return err
	}

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
		cb(self.Handlers(b))
	}

	return nil
}

func (self *NeighborAPI) Handlers(b []byte) error {
	msg := strings.Split(string(b), " ")

	log.Println("received: " + string(b))

	if msg[0] == "scrooge_hello" {
		return self.HelloHandler(msg)
	} else {
		return errors.New("unrecognized message type")
	}
}

func (self *NeighborAPI) HelloHandler(msg []string) error {
	helloMessage, err := serialization.ParseHello(msg)
	if err != nil {
		return err
	}

	neighbor := self.Neighbors[helloMessage.PublicKey]
	if neighbor == nil {
		neighbor = &types.Neighbor{
			PublicKey: helloMessage.PublicKey,
		}
	}

	if neighbor.Seqnum > helloMessage.Seqnum {
		return errors.New("sequence number too low")
	}

	neighbor.Seqnum = helloMessage.Seqnum
	neighbor.ControlAddress = helloMessage.ControlAddress

	addr, err := net.ResolveUDPAddr("udp6", neighbor.ControlAddress)
	if err != nil {
		return err
	}

	err = self.Network.SendUDP(addr, serialization.FmtHello(self.Account))
	if err != nil {
		return err
	}
	return nil
}

// McastHello sends an scrooge_hello packet to the multicast UDP address on a given interface.
func (self *NeighborAPI) McastHello(
	mCastPort int,
	iface *net.Interface,
	cb func(*types.Neighbor, error),
) error {
	err := self.Network.SendUDP(&net.UDPAddr{
		IP:   net.ParseIP("ff02::1"),
		Port: mCastPort,
		Zone: iface.Name,
	}, serialization.FmtHello(self.Account))
	if err != nil {
		return err
	}
	return nil
}

// func (self *NeighborAPI) sendUDP(
// 	addr *net.UDPAddr,
// 	s string,
// ) error {
// 	conn, err := net.DialUDP(
// 		"udp6",
// 		nil,
// 		addr,
// 	)
// 	defer conn.Close()
// 	if err != nil {
// 		return err
// 	}

// 	conn.Write([]byte(s))
// 	log.Println("sent: " + s)
// 	return nil
// }
