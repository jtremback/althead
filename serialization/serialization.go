package serialization

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"strconv"

	"github.com/agl/ed25519"
	"github.com/incentivized-mesh-infrastructure/scrooge/types"
)

// scrooge_hello[_confirm] <publicKey> <control address> <seqnum> <signature>
func FmtHello(
	msg types.HelloMessage,
	privateKey [ed25519.PrivateKeySize]byte,
) (string, error) {
	var msgType string

	if msg.Confirm {
		msgType = "scrooge_hello_confirm"
	} else {
		msgType = "scrooge_hello"
	}

	s := fmt.Sprintf(
		"%v %v %v %v",
		msgType,
		base64.StdEncoding.EncodeToString(msg.PublicKey[:]),
		msg.ControlAddress.String(),
		msg.Seqnum,
	)

	sig := ed25519.Sign(&privateKey, []byte(s))

	return s + " " + base64.StdEncoding.EncodeToString(sig[:]), nil
}

func ParseHello(msg []string) (*types.HelloMessage, error) {
	var confirm bool
	if msg[0] == "scrooge_hello" {
		confirm = false
	} else if msg[0] == "scrooge_hello_confirm" {
		confirm = true
	} else {
		return nil, errors.New("Not a scrooge_hello or scrooge_hello_confirm message")
	}

	messageMetadata, err := verifyMessage(msg)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp6", msg[2])
	if err != nil {
		return nil, err
	}

	h := &types.HelloMessage{
		MessageMetadata: *messageMetadata,
		ControlAddress:  *addr,
		Confirm:         confirm,
	}

	log.Printf("parsed HelloMessage: %+v\n", h)

	return h, nil
}

// scrooge_tunnel[_confirm] <publicKey> <tunnel publicKey> <tunnel endpoint> <seq num> <signature>
func FmtTunnel(
	msg types.TunnelMessage,
	privateKey [ed25519.PrivateKeySize]byte,
) (string, error) {
	var msgType string

	if msg.Confirm {
		msgType = "scrooge_tunnel_confirm"
	} else {
		msgType = "scrooge_tunnel"
	}

	s := fmt.Sprintf(
		"%v %v %v %v %v",
		msgType,
		base64.StdEncoding.EncodeToString(msg.PublicKey[:]),
		msg.TunnelPublicKey,
		msg.TunnelEndpoint,
		msg.Seqnum,
	)

	sig := ed25519.Sign(&privateKey, []byte(s))

	return s + " " + base64.StdEncoding.EncodeToString(sig[:]), nil
}

func ParseTunnel(msg []string) (*types.TunnelMessage, error) {
	var confirm bool
	if msg[0] == "scrooge_tunnel" {
		confirm = false
	} else if msg[0] == "scrooge_tunnel_confirm" {
		confirm = true
	} else {
		return nil, errors.New("Not a scrooge_tunnel or scrooge_tunnel_confirm message")
	}

	messageMetadata, err := verifyMessage(msg)
	if err != nil {
		return nil, err
	}

	m := &types.TunnelMessage{
		MessageMetadata: *messageMetadata,
		TunnelPublicKey: msg[2],
		TunnelEndpoint:  msg[3],
		Confirm:         confirm,
	}

	log.Printf("parsed TunnelMessage: %+v\n", m)

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

	msgWithOutSig := strings.Join(msg[:len(msg)-1], " ")

	if !ed25519.Verify(&publicKey, []byte(msgWithOutSig), &signature) {
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
