package findPeersMCast

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/jtremback/althea/ed25519-wrapper"
	"github.com/jtremback/althea/types"
)

func Listen(
	iface *net.Interface,
	mcastPort int,
	account types.Account,
	cb func(*types.Peer, error),
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
		peer, err := readUDPHello(conn)
		if err != nil {
			cb(nil, err)
			continue
		}
		cb(peer, nil)
	}

	return nil
}

func readUDPHello(conn *net.UDPConn) (*types.Peer, error) {
	b := make([]byte, 64)
	_, addr, err := conn.ReadFromUDP(b)
	if err != nil {
		return nil, err
	}

	msg := strings.Split(string(b), " ")

	log.Println("received: "+string(b), "from: ", addr)

	if msg[0] == "althea_hello" {
		peer, err := parseHello(msg)
		if err != nil {
			return nil, err
		}
		return peer, nil
	} else {
		return nil, errors.New("unrecognized message type")
	}
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

// althea_hello <control address> <control pubkey> <tunnel address> <tunnel pubkey> <signature>
func printHello(account types.Account) string {
	sig := ed25519.Sign(account.ControlPrivkey, concatByteSlices(
		[]byte("althea_hello"),
		[]byte(account.ControlAddress),
		account.ControlPubkey,
		[]byte(account.TunnelAddress),
		account.TunnelPubkey,
	))

	return fmt.Sprintf(
		"althea_hello %v %v %v %v %v",
		account.ControlAddress,
		base64.StdEncoding.EncodeToString(account.ControlPubkey),
		account.TunnelAddress,
		base64.StdEncoding.EncodeToString(account.TunnelPubkey),
		base64.StdEncoding.EncodeToString(sig),
	)
}

func concatByteSlices(slices ...[]byte) []byte {
	var slice []byte
	for _, s := range slices {
		slice = append(slice, s...)
	}
	return slice
}

func parseHello(msg []string) (*types.Peer, error) {
	controlPubkey, err := base64.StdEncoding.DecodeString(msg[2])
	if err != nil {
		return nil, err
	}

	tunnelPubkey, err := base64.StdEncoding.DecodeString(msg[4])
	if err != nil {
		return nil, err
	}

	sig, err := base64.StdEncoding.DecodeString(msg[5])
	if err != nil {
		return nil, err
	}

	peer := &types.Peer{
		ControlAddress: msg[1],
		ControlPubkey:  controlPubkey,
		TunnelAddress:  msg[3],
		TunnelPubkey:   tunnelPubkey,
	}

	if !ed25519.Verify(controlPubkey, concatByteSlices(
		[]byte("althea_hello"),
		[]byte(peer.ControlAddress),
		peer.ControlPubkey,
		[]byte(peer.TunnelAddress),
		[]byte(peer.TunnelPubkey),
	), sig) {
		return nil, errors.New("signature not valid")
	}

	return peer, nil
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
	cb func(*types.Peer, error),
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
		peer, err := readUDPHello(conn)
		if err != nil {
			cb(nil, err)
			continue
		}
		cb(peer, nil)
	}
}
