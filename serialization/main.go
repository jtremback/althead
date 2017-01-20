package serialization

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"strconv"

	"github.com/agl/ed25519"
	"github.com/jtremback/scrooge/types"
)

// scrooge_hello <publicKey> <control address> <seqnum> <signature>
func FmtHello(account *types.Account) string {
	msg := fmt.Sprintf(
		"scrooge_hello %v %v %v",
		base64.StdEncoding.EncodeToString(account.PublicKey[:]),
		account.ControlAddress,
		account.Seqnum,
	)

	sig := ed25519.Sign(&account.PrivateKey, []byte(msg))

	return msg + " " + base64.StdEncoding.EncodeToString(sig[:])
}

func FmtIHU(account *types.Account) string {
	msg := fmt.Sprintf(
		"scrooge_ihu %v %v %v",
		base64.StdEncoding.EncodeToString(account.PublicKey[:]),
		account.ControlAddress,
		account.Seqnum,
	)

	sig := ed25519.Sign(&account.PrivateKey, []byte(msg))

	return msg + " " + base64.StdEncoding.EncodeToString(sig[:])
}

func ParseHello(msg []string) (*types.HelloMessage, error) {
	pk, err := base64.StdEncoding.DecodeString(msg[1])
	if err != nil {
		return nil, err
	}

	s, err := base64.StdEncoding.DecodeString(msg[4])
	if err != nil {
		return nil, err
	}

	publicKey := types.BytesToPublicKey(pk)
	sig := types.BytesToSignature(s)

	if !ed25519.Verify(&publicKey, []byte(strings.Join(msg[:3], " ")), &sig) {
		return nil, errors.New("signature not valid")
	}

	controlAddress := msg[2]

	seqnum, err := strconv.ParseUint(msg[3], 10, 64)
	if err != nil {
		return nil, err
	}

	return &types.HelloMessage{
		PublicKey:      publicKey,
		Seqnum:         seqnum,
		ControlAddress: controlAddress,
		Signature:      sig,
	}, nil
}
