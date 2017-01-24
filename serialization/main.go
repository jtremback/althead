package serialization

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"strconv"

	"net"

	"github.com/agl/ed25519"
	"github.com/jtremback/scrooge/types"
)

// scrooge_hello <publicKey> <control address> <seqnum> <signature>
func FmtHello(
	account *types.Account,
	iface *net.Interface,
) string {
	msg := types.HelloMessage{
		MessageMetadata: types.MessageMetadata{
			Seqnum:    account.Seqnum,
			PublicKey: account.PublicKey,
		},
		ControlAddress: account.ControlAddresses[iface.Name],
	}

	s := fmt.Sprintf(
		"scrooge_hello %v %v %v",
		base64.StdEncoding.EncodeToString(msg.PublicKey[:]),
		msg.ControlAddress,
		msg.Seqnum,
	)

	sig := ed25519.Sign(&account.PrivateKey, []byte(s))

	return s + " " + base64.StdEncoding.EncodeToString(sig[:])
}

func ParseHello(msg []string) (*types.HelloMessage, error) {
	messageMetadata, err := verifyMessage(msg)
	if err != nil {
		return nil, err
	}

	h := &types.HelloMessage{
		MessageMetadata: *messageMetadata,
		ControlAddress:  msg[2],
	}

	return h, nil
}

// scrooge_hello_confirm <publicKey> <control address> <seqnum> <signature>
func FmtHelloConfirm(
	account *types.Account,
	iface *net.Interface,
) string {
	msg := types.HelloConfirmMessage{
		MessageMetadata: types.MessageMetadata{
			Seqnum:    account.Seqnum,
			PublicKey: account.PublicKey,
		},
		ControlAddress: account.ControlAddresses[iface.Name],
	}

	s := fmt.Sprintf(
		"scrooge_hello_confirm %v %v %v",
		base64.StdEncoding.EncodeToString(msg.PublicKey[:]),
		msg.ControlAddress,
		msg.Seqnum,
	)

	sig := ed25519.Sign(&account.PrivateKey, []byte(s))

	return s + " " + base64.StdEncoding.EncodeToString(sig[:])
}

func ParseHelloConfirm(msg []string) (*types.HelloConfirmMessage, error) {
	messageMetadata, err := verifyMessage(msg)
	if err != nil {
		return nil, err
	}

	h := &types.HelloConfirmMessage{
		MessageMetadata: *messageMetadata,
		ControlAddress:  msg[2],
	}

	return h, nil
}

// scrooge_tunnel <publicKey> <tunnel publicKey> <tunnel endpoint> <seq num> <signature>
func FmtTunnel(
	account *types.Account,
	neighbor *types.Neighbor,
	iface *net.Interface,
) string {
	m := types.TunnelMessage{
		MessageMetadata: types.MessageMetadata{
			PublicKey: account.PublicKey,
			Seqnum:    account.Seqnum,
		},
		TunnelEndpoint:  neighbor.Tunnel.Endpoint,
		TunnelPublicKey: neighbor.Tunnel.PublicKey,
	}

	msg := fmt.Sprintf(
		"scrooge_tunnel %v %v %v",
		base64.StdEncoding.EncodeToString(m.PublicKey[:]),
		m.TunnelPublicKey,
		m.TunnelEndpoint,
		m.Seqnum,
	)

	sig := ed25519.Sign(&account.PrivateKey, []byte(msg))

	return msg + " " + base64.StdEncoding.EncodeToString(sig[:])
}

func ParseTunnel(msg []string) (*types.TunnelMessage, error) {
	messageMetadata, err := verifyMessage(msg)
	if err != nil {
		return nil, err
	}

	m := &types.TunnelMessage{
		MessageMetadata: *messageMetadata,
		TunnelPublicKey: msg[2],
		TunnelEndpoint:  msg[3],
	}

	return m, nil
}

// scrooge_tunnel_confirm <publicKey> <tunnel publicKey> <tunnel endpoint> <seq num> <signature>
func FmtTunnelConfirm(
	account *types.Account,
	neighbor *types.Neighbor,
	iface *net.Interface,
) string {
	m := types.TunnelConfirmMessage{
		MessageMetadata: types.MessageMetadata{
			PublicKey: account.PublicKey,
			Seqnum:    account.Seqnum,
		},
		TunnelEndpoint:  neighbor.Tunnel.Endpoint,
		TunnelPublicKey: neighbor.Tunnel.PublicKey,
	}

	msg := fmt.Sprintf(
		"scrooge_tunnel_confirm %v %v %v",
		base64.StdEncoding.EncodeToString(m.PublicKey[:]),
		m.TunnelPublicKey,
		m.TunnelEndpoint,
		m.Seqnum,
	)

	sig := ed25519.Sign(&account.PrivateKey, []byte(msg))

	return msg + " " + base64.StdEncoding.EncodeToString(sig[:])
}

func ParseTunnelConfirm(msg []string) (*types.TunnelConfirmMessage, error) {
	messageMetadata, err := verifyMessage(msg)
	if err != nil {
		return nil, err
	}

	m := &types.TunnelConfirmMessage{
		MessageMetadata: *messageMetadata,
		TunnelPublicKey: msg[2],
		TunnelEndpoint:  msg[3],
	}

	return m, nil
}

func verifyMessage(msg []string) (*types.MessageMetadata, error) {
	pk, err := base64.StdEncoding.DecodeString(msg[1])
	if err != nil {
		return nil, err
	}

	sig, err := base64.StdEncoding.DecodeString(msg[len(msg)-1])
	if err != nil {
		return nil, err
	}

	publicKey := types.BytesToPublicKey(pk)
	signature := types.BytesToSignature(sig)

	if !ed25519.Verify(&publicKey, []byte(strings.Join(msg[:len(msg)-1], " ")), &signature) {
		return nil, errors.New("signature not valid")
	}

	seqnum, err := strconv.ParseUint(msg[len(msg)-2], 10, 64)
	if err != nil {
		return nil, err
	}

	messageMetadata := types.MessageMetadata{
		PublicKey: publicKey,
		Seqnum:    seqnum,
		Signature: signature,
	}

	return &messageMetadata, nil
}
