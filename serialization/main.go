package serialization

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/jtremback/althea/ed25519-wrapper"
	"github.com/jtremback/althea/types"
)

func concatByteSlices(slices ...[]byte) []byte {
	var slice []byte
	for _, s := range slices {
		slice = append(slice, s...)
	}
	return slice
}

// althea_hello <control address> <control pubkey> <tunnel address> <tunnel pubkey> <signature>
func FmtHello(account types.Account) string {
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

func ParseHello(msg []string) (*types.Neighbor, error) {
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

	neighbor := &types.Neighbor{
		ControlAddress: msg[1],
		ControlPubkey:  controlPubkey,
		TunnelAddress:  msg[3],
		TunnelPubkey:   tunnelPubkey,
	}

	if !ed25519.Verify(controlPubkey, concatByteSlices(
		[]byte("althea_hello"),
		[]byte(neighbor.ControlAddress),
		neighbor.ControlPubkey,
		[]byte(neighbor.TunnelAddress),
		[]byte(neighbor.TunnelPubkey),
	), sig) {
		return nil, errors.New("signature not valid")
	}

	return neighbor, nil
}
