package findPeersMCast

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/agl/ed25519"
	"github.com/jtremback/althea/types"
)

func Listen(
	iface *net.Interface,
	mcastPort int,
	account types.Account,
	cb func(types.Peer, error),
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
			sendUDP(addr, printHello(account))
		}
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

func printHello(account types.Account) string {
	// althea_hello <control pubkey> <control address> <tunnel pubkey> <tunnel address> <signature>
	unsigned := fmt.Sprintf(
		"althea_hello %v %v %v %v ",
		base64.URLEncoding.EncodeToString(account.ControlPubkey[:]),
		account.TunnelAddress,
		base64.URLEncoding.EncodeToString(ed25519.Sign(&account.ControlPrivkey, []byte(account.TunnelAddress))[:]),
	)

	sig := base64.URLEncoding.EncodeToString(ed25519.Sign(&account.ControlPrivkey, []byte(unsigned))[:])

	return unsigned + sig
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
