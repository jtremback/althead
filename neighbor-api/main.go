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
		SendMulticastUDP(*net.Interface, int, string) error
	}
}

func (self *NeighborAPI) Handlers(b []byte, iface *net.Interface) error {
	msg := strings.Split(string(b), " ")

	log.Println("received: " + string(b))

	if msg[0] == "scrooge_hello" {
		return self.helloHandler(msg, iface)
	} else {
		return errors.New("unrecognized message type")
	}
}

func (self *NeighborAPI) helloHandler(
	msg []string,
	iface *net.Interface,
) error {
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

	err = self.SendHelloConfirm(addr, iface)
	if err != nil {
		return err
	}
	return nil
}

// func (self *NeighborAPI) helloHandler(
// 	msg []string,
// 	iface *net.Interface,
// ) error {
// 	helloMessage, err := serialization.ParseHello(msg)
// 	if err != nil {
// 		return err
// 	}

// 	neighbor := self.Neighbors[helloMessage.PublicKey]
// 	if neighbor == nil {
// 		neighbor = &types.Neighbor{
// 			PublicKey: helloMessage.PublicKey,
// 		}
// 	}

// 	if neighbor.Seqnum > helloMessage.Seqnum {
// 		return errors.New("sequence number too low")
// 	}

// 	neighbor.Seqnum = helloMessage.Seqnum
// 	neighbor.ControlAddress = helloMessage.ControlAddress

// 	addr, err := net.ResolveUDPAddr("udp6", neighbor.ControlAddress)
// 	if err != nil {
// 		return err
// 	}

// 	// if helloMessage.Confirm != types.{
// 	// 	err = self.SendHelloConfirm(addr, iface)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// }
// 	return nil
// }

func (self *NeighborAPI) SendHelloConfirm(
	addr *net.UDPAddr,
	iface *net.Interface,
) error {
	err := self.Network.SendUDP(addr, serialization.FmtHelloConfirm(self.Account, iface))
	if err != nil {
		return err
	}
	return nil
}

func (self *NeighborAPI) SendHello(
	addr *net.UDPAddr,
	iface *net.Interface,
) error {
	err := self.Network.SendUDP(addr, serialization.FmtHello(self.Account, iface))
	if err != nil {
		return err
	}
	return nil
}

func (self *NeighborAPI) SendMcastHello(
	iface *net.Interface,
	port int,
) error {
	err := self.Network.SendMulticastUDP(iface, port, serialization.FmtHello(self.Account, iface))
	if err != nil {
		return err
	}
	return nil
}
