package serialization

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"strconv"

	"github.com/agl/ed25519"
	"github.com/incentivized-mesh-infrastructure/scrooge/types"
)

// scrooge_hello[_confirm] <SourcePublicKey> <seqnum> <signature>
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
		"%v %v %v",
		msgType,
		base64.StdEncoding.EncodeToString(msg.SourcePublicKey[:]),
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

	messageMetadata, err := verifyMessage(msg, true)
	if err != nil {
		return nil, err
	}

	h := &types.HelloMessage{
		MessageMetadata: *messageMetadata,
		Confirm:         confirm,
	}

	log.Printf("parsed HelloMessage: %+v\n", h)

	return h, nil
}

// scrooge_tunnel[_confirm] <sourcePublicKey> <destinationPublicKey> <tunnel publicKey> <tunnel endpoint> <seq num> <signature>
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
		"%v %v %v %v %v %v",
		msgType,
		base64.StdEncoding.EncodeToString(msg.SourcePublicKey[:]),
		base64.StdEncoding.EncodeToString(msg.DestinationPublicKey[:]),
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

	messageMetadata, err := verifyMessage(msg, false)
	if err != nil {
		return nil, err
	}

	m := &types.TunnelMessage{
		MessageMetadata: *messageMetadata,
		TunnelPublicKey: msg[3],
		TunnelEndpoint:  msg[4],
		Confirm:         confirm,
	}

	log.Printf("parsed TunnelMessage: %+v\n", m)

	return m, nil
}

func verifyMessage(msg []string, broadcast bool) (*types.MessageMetadata, error) {
	sig, err := base64.StdEncoding.DecodeString(msg[len(msg)-1])
	if err != nil {
		return nil, err
	}

	spk, err := base64.StdEncoding.DecodeString(msg[1])
	if err != nil {
		return nil, err
	}

	var destinationPublicKey [ed25519.PublicKeySize]byte

	if !broadcast {
		dpk, err := base64.StdEncoding.DecodeString(msg[1])
		if err != nil {
			return nil, err
		}
		destinationPublicKey = types.BytesToPublicKey(dpk)
	}

	sourcePublicKey := types.BytesToPublicKey(spk)
	signature := types.BytesToSignature(sig)

	msgWithOutSig := strings.Join(msg[:len(msg)-1], " ")

	if !ed25519.Verify(&sourcePublicKey, []byte(msgWithOutSig), &signature) {
		return nil, errors.New("signature not valid")
	}

	seqnum, err := strconv.ParseUint(msg[len(msg)-2], 10, 64)
	if err != nil {
		return nil, err
	}

	messageMetadata := types.MessageMetadata{
		SourcePublicKey:      sourcePublicKey,
		DestinationPublicKey: destinationPublicKey,
		Seqnum:               seqnum,
		Signature:            signature,
	}

	return &messageMetadata, nil
}
