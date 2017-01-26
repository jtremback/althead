package serialization

import (
	"testing"

	"strings"

	"github.com/agl/ed25519"
	"github.com/jtremback/scrooge/types"
)

var (
	pubkey1                     = &[ed25519.PublicKeySize]byte{44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	privkey1                    = &[ed25519.PrivateKeySize]byte{112, 69, 149, 144, 72, 233, 25, 188, 124, 215, 67, 200, 213, 237, 133, 127, 215, 253, 230, 134, 26, 202, 25, 214, 36, 19, 233, 87, 212, 169, 119, 226, 44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	pubkey2                     = &[ed25519.PublicKeySize]byte{175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123, 162}
	privkey2                    = &[ed25519.PrivateKeySize]byte{13, 170, 251, 93, 50, 201, 207, 72, 224, 172, 35, 48, 16, 245, 116, 20, 88, 33, 155, 12, 226, 126, 59, 36, 184, 111, 95, 87, 156, 104, 140, 243, 175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123}
	helloMessage                = "scrooge_hello LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= 1.1.1.1 12 MGV9+pZfqUE9vVUFBYmV9plPbXYXan7yIkIt3sF0Zvuz+rEO+rLyqcWobCPRdjRIab+bZpIc+nPp8MWjNWa2CA=="
	helloConfirmMessage         = "scrooge_hello_confirm LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= 1.1.1.1 12 HFIZJXcDnCIrIUhAst5DSffBgd7b3LZyA8ymQu1iPwKbIqM7FC+vt0js1GbpCZHBf0Kk5tEHb8hsWTKGufThBg=="
	tunnelMessage               = "scrooge_tunnel LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= derp 2.2.2.2 12 vA7wEuO/7h92hBv8ZVQq9KjtshJrsebpWCU7BGnFuuon1XsH2FWOABQYU/aBmn+gIfTJUv9Pu8khUAKAVLzNDA=="
	tunnelConfirmMessage        = "scrooge_tunnel_confirm LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= derp 2.2.2.2 12 GJ90TCQj1KFJz6Y369orO5Y5I4HHqhJPm4Ew7/BU8BYRnWqtG2FA/1TlqgwwwMGZDkN7RU/40sOxnmYf26bgCQ=="
	iface1                      = "eth0"
	controlAddress1             = "1.1.1.1"
	seqnum1              uint64 = 12
	tunnelEndpoint1             = "2.2.2.2"
	tunnelPubkey1               = "derp"
)

func TestFmtHello(t *testing.T) {
	testFmtHello(t, false)
}

func TestFmtHelloConfirm(t *testing.T) {
	testFmtHello(t, true)
}

func testFmtHello(t *testing.T, confirm bool) {
	acct := &types.Account{
		Seqnum:     seqnum1,
		PublicKey:  *pubkey1,
		PrivateKey: *privkey1,
	}

	msg, err := FmtHello(acct, controlAddress1, confirm)
	if err != nil {
		t.Error(err)
	}

	var realMsg string

	if confirm {
		realMsg = helloConfirmMessage
	} else {
		realMsg = helloMessage
	}

	if msg != realMsg {
		t.Error("Message format incorrect: " + msg)
	}
}

func TestParseHello(t *testing.T) {
	testParseHello(t, false)
}

func TestParseHelloConfirm(t *testing.T) {
	testParseHello(t, true)
}

func testParseHello(t *testing.T, confirm bool) {
	var realMsg string

	if confirm {
		realMsg = helloConfirmMessage
	} else {
		realMsg = helloMessage
	}

	msg, err := ParseHello(strings.Split(realMsg, " "))
	if err != nil {
		t.Fatal(err)
	}
	if msg.PublicKey != *pubkey1 {
		t.Error("msg.PublicKey incorrect")
	}
	if msg.ControlAddress != controlAddress1 {
		t.Error("msg.ControlAddress incorrect")
	}
	if msg.Seqnum != seqnum1 {
		t.Error("msg.Seqnum incorrect")
	}

	var sig [ed25519.SignatureSize]byte

	if confirm {
		sig = [ed25519.SignatureSize]byte{28, 82, 25, 37, 119, 3, 156, 34, 43, 33, 72, 64, 178, 222, 67, 73, 247, 193, 129, 222, 219, 220, 182, 114, 3, 204, 166, 66, 237, 98, 63, 2, 155, 34, 163, 59, 20, 47, 175, 183, 72, 236, 212, 102, 233, 9, 145, 193, 127, 66, 164, 230, 209, 7, 111, 200, 108, 89, 50, 134, 185, 244, 225, 6}
	} else {
		sig = [ed25519.SignatureSize]byte{48, 101, 125, 250, 150, 95, 169, 65, 61, 189, 85, 5, 5, 137, 149, 246, 153, 79, 109, 118, 23, 106, 126, 242, 34, 66, 45, 222, 193, 116, 102, 251, 179, 250, 177, 14, 250, 178, 242, 169, 197, 168, 108, 35, 209, 118, 52, 72, 105, 191, 155, 102, 146, 28, 250, 115, 233, 240, 197, 163, 53, 102, 182, 8}
	}

	if msg.Signature != sig {
		t.Error("msg.Signature incorrect: ", msg.Signature, sig)
	}
}

func TestFmtTunnel(t *testing.T) {
	testFmtTunnel(t, false)
}

func TestFmtTunnelConfirm(t *testing.T) {
	testFmtTunnel(t, true)
}

func testFmtTunnel(t *testing.T, confirm bool) {
	acct := &types.Account{
		Seqnum:     seqnum1,
		PublicKey:  *pubkey1,
		PrivateKey: *privkey1,
	}

	neigh := &types.Neighbor{
		Seqnum:    seqnum1,
		PublicKey: *pubkey1,
	}

	neigh.Tunnel.Endpoint = tunnelEndpoint1
	neigh.Tunnel.PublicKey = tunnelPubkey1

	msg, err := FmtTunnel(acct, neigh, confirm)
	if err != nil {
		t.Error(err)
	}

	var realMsg string

	if confirm {
		realMsg = tunnelConfirmMessage
	} else {
		realMsg = tunnelMessage
	}

	if msg != realMsg {
		t.Error("Message format incorrect: " + msg)
	}
}

func TestParseTunnel(t *testing.T) {
	testParseTunnel(t, false)
}

func TestParseTunnelConfirm(t *testing.T) {
	testParseTunnel(t, true)
}

func testParseTunnel(t *testing.T, confirm bool) {
	var realMsg string

	if confirm {
		realMsg = tunnelConfirmMessage
	} else {
		realMsg = tunnelMessage
	}

	msg, err := ParseTunnel(strings.Split(realMsg, " "))
	if err != nil {
		t.Fatal(err)
	}
	if msg.PublicKey != *pubkey1 {
		t.Error("msg.PublicKey incorrect")
	}
	if msg.TunnelEndpoint != tunnelEndpoint1 {
		t.Error("msg.TunnelEndpoint incorrect", msg.TunnelEndpoint)
	}
	if msg.TunnelPublicKey != tunnelPubkey1 {
		t.Error("msg.TunnelPublicKey incorrect", msg.TunnelPublicKey)
	}
	if msg.Seqnum != seqnum1 {
		t.Error("msg.Seqnum incorrect")
	}

	var sig [ed25519.SignatureSize]byte

	if confirm {
		sig = [ed25519.SignatureSize]byte{24, 159, 116, 76, 36, 35, 212, 161, 73, 207, 166, 55, 235, 218, 43, 59, 150, 57, 35, 129, 199, 170, 18, 79, 155, 129, 48, 239, 240, 84, 240, 22, 17, 157, 106, 173, 27, 97, 64, 255, 84, 229, 170, 12, 48, 192, 193, 153, 14, 67, 123, 69, 79, 248, 210, 195, 177, 158, 102, 31, 219, 166, 224, 9}
	} else {
		sig = [ed25519.SignatureSize]byte{188, 14, 240, 18, 227, 191, 238, 31, 118, 132, 27, 252, 101, 84, 42, 244, 168, 237, 178, 18, 107, 177, 230, 233, 88, 37, 59, 4, 105, 197, 186, 234, 39, 213, 123, 7, 216, 85, 142, 0, 20, 24, 83, 246, 129, 154, 127, 160, 33, 244, 201, 82, 255, 79, 187, 201, 33, 80, 2, 128, 84, 188, 205, 12}
	}

	if msg.Signature != sig {
		t.Error("msg.Signature incorrect: ", msg.Signature, sig)
	}
}
