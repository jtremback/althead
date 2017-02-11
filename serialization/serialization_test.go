package serialization

import (
	"testing"

	"strings"

	"net"

	"github.com/agl/ed25519"
	"github.com/incentivized-mesh-infrastructure/scrooge/types"
)

var (
	pubkey1              = &[ed25519.PublicKeySize]byte{44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	privkey1             = &[ed25519.PrivateKeySize]byte{112, 69, 149, 144, 72, 233, 25, 188, 124, 215, 67, 200, 213, 237, 133, 127, 215, 253, 230, 134, 26, 202, 25, 214, 36, 19, 233, 87, 212, 169, 119, 226, 44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	pubkey2              = &[ed25519.PublicKeySize]byte{175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123, 162}
	privkey2             = &[ed25519.PrivateKeySize]byte{13, 170, 251, 93, 50, 201, 207, 72, 224, 172, 35, 48, 16, 245, 116, 20, 88, 33, 155, 12, 226, 126, 59, 36, 184, 111, 95, 87, 156, 104, 140, 243, 175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123}
	helloMessage         = "scrooge_hello LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= 1.1.1.1:8000 12 bLKZBGYmwOiKnFflOEuZ0iRmIA0YtPeAyJM/wyh10pPeH+gYG2tMdDDCdhwYQKJtB3CTZ+Mn9s5IfxnbQB7jCQ=="
	helloConfirmMessage  = "scrooge_hello_confirm LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= 1.1.1.1:8000 12 WV4qgi9H5s5mbXg5B1LvKHgXepm+HkZmymqavQza7+1CvTLb+6R3yAJ1GMFPCetPfX8GLNTS4AJiN2tBhf7WAg=="
	tunnelMessage        = "scrooge_tunnel LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= derp 2.2.2.2:8000 12 6ylkg9zJjk0yYvudq2mUQVi+CaFo9faLsT8Yjiy/g95mTzbrul3eUJabrdiRm/gnbyS5lQHff5TjI2DMfaOZBw=="
	tunnelConfirmMessage = "scrooge_tunnel_confirm LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= derp 2.2.2.2:8000 12 DEXqnvdy8ppaulkXO7ivJr+bQphGiSL3zAGFuMgvg4wvbuUlcAX+3+6bYE2xwssUmjtCyK7Z+XkKGhzHC5ChBw=="
	iface1               = "eth0"
	controlAddress1      = net.UDPAddr{
		IP:   net.ParseIP("1.1.1.1"),
		Port: 8000,
	}
	seqnum1         uint64 = 12
	tunnelEndpoint1        = "2.2.2.2:8000"
	tunnelPubkey1          = "derp"
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

	msg := types.HelloMessage{
		MessageMetadata: types.MessageMetadata{
			Seqnum:    acct.Seqnum,
			PublicKey: acct.PublicKey,
		},
		ControlAddress: controlAddress1,
		Confirm:        confirm,
	}

	s, err := FmtHello(msg, acct.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	var realMsg string

	if confirm {
		realMsg = helloConfirmMessage
	} else {
		realMsg = helloMessage
	}

	if s != realMsg {
		t.Fatal("Message format incorrect: " + s)
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
		t.Fatal("msg.PublicKey incorrect")
	}
	if msg.ControlAddress.String() != controlAddress1.String() {
		t.Fatal("msg.ControlAddress incorrect", msg.ControlAddress.String(), controlAddress1.String())
	}
	if msg.Seqnum != seqnum1 {
		t.Fatal("msg.Seqnum incorrect")
	}

	var sig [ed25519.SignatureSize]byte

	if confirm {
		sig = [ed25519.SignatureSize]byte{89, 94, 42, 130, 47, 71, 230, 206, 102, 109, 120, 57, 7, 82, 239, 40, 120, 23, 122, 153, 190, 30, 70, 102, 202, 106, 154, 189, 12, 218, 239, 237, 66, 189, 50, 219, 251, 164, 119, 200, 2, 117, 24, 193, 79, 9, 235, 79, 125, 127, 6, 44, 212, 210, 224, 2, 98, 55, 107, 65, 133, 254, 214, 2}
	} else {
		sig = [ed25519.SignatureSize]byte{108, 178, 153, 4, 102, 38, 192, 232, 138, 156, 87, 229, 56, 75, 153, 210, 36, 102, 32, 13, 24, 180, 247, 128, 200, 147, 63, 195, 40, 117, 210, 147, 222, 31, 232, 24, 27, 107, 76, 116, 48, 194, 118, 28, 24, 64, 162, 109, 7, 112, 147, 103, 227, 39, 246, 206, 72, 127, 25, 219, 64, 30, 227, 9}
	}

	if msg.Signature != sig {
		t.Fatal("msg.Signature incorrect: ", msg.Signature, sig)
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

	neighbor := &types.Neighbor{
		Seqnum:    seqnum1,
		PublicKey: *pubkey1,
	}

	neighbor.Tunnel.Endpoint = tunnelEndpoint1
	neighbor.Tunnel.PublicKey = tunnelPubkey1

	msg := types.TunnelMessage{
		MessageMetadata: types.MessageMetadata{
			PublicKey: acct.PublicKey,
			Seqnum:    acct.Seqnum,
		},
		TunnelEndpoint:  neighbor.Tunnel.Endpoint,
		TunnelPublicKey: neighbor.Tunnel.PublicKey,
		Confirm:         confirm,
	}

	s, err := FmtTunnel(msg, acct.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	var realMsg string

	if confirm {
		realMsg = tunnelConfirmMessage
	} else {
		realMsg = tunnelMessage
	}

	if s != realMsg {
		t.Fatal("Message format incorrect: " + s)
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
		t.Fatal("msg.PublicKey incorrect")
	}
	if msg.TunnelEndpoint != tunnelEndpoint1 {
		t.Fatal("msg.TunnelEndpoint incorrect", msg.TunnelEndpoint)
	}
	if msg.TunnelPublicKey != tunnelPubkey1 {
		t.Fatal("msg.TunnelPublicKey incorrect", msg.TunnelPublicKey)
	}
	if msg.Seqnum != seqnum1 {
		t.Fatal("msg.Seqnum incorrect")
	}

	var sig [ed25519.SignatureSize]byte

	if confirm {
		sig = [ed25519.SignatureSize]byte{12, 69, 234, 158, 247, 114, 242, 154, 90, 186, 89, 23, 59, 184, 175, 38, 191, 155, 66, 152, 70, 137, 34, 247, 204, 1, 133, 184, 200, 47, 131, 140, 47, 110, 229, 37, 112, 5, 254, 223, 238, 155, 96, 77, 177, 194, 203, 20, 154, 59, 66, 200, 174, 217, 249, 121, 10, 26, 28, 199, 11, 144, 161, 7}
	} else {
		sig = [ed25519.SignatureSize]byte{235, 41, 100, 131, 220, 201, 142, 77, 50, 98, 251, 157, 171, 105, 148, 65, 88, 190, 9, 161, 104, 245, 246, 139, 177, 63, 24, 142, 44, 191, 131, 222, 102, 79, 54, 235, 186, 93, 222, 80, 150, 155, 173, 216, 145, 155, 248, 39, 111, 36, 185, 149, 1, 223, 127, 148, 227, 35, 96, 204, 125, 163, 153, 7}
	}

	if msg.Signature != sig {
		t.Fatal("msg.Signature incorrect: ", msg.Signature, sig)
	}
}
