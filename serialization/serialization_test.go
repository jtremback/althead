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
	tunnelMessage        = "scrooge_tunnel LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= r24MX1Kp720po7ddTcUjKSPLXsjYBimBqgwIYdMce6I= flerp 3.3.3.3:8000 12 kaklnXMRGqs/9GCh8D0TWlYRjIMXNSugJQYuHolfa0ERsHbLpstCUs6kzXISwb6TX2E/ultWHjkDd1vkK6BrBA=="
	tunnelConfirmMessage = "scrooge_tunnel_confirm LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= r24MX1Kp720po7ddTcUjKSPLXsjYBimBqgwIYdMce6I= flerp 3.3.3.3:8000 12 nKOU7/fQLFgWffp6VQ0H2WKz68L8NFXy7CW9ku8zy7BJqGpmTnjaeJiomBkPpwo/6nmPIL69sAIW8iOv6vqvDw=="
	iface1               = "eth0"
	controlAddress1      = net.UDPAddr{
		IP:   net.ParseIP("1.1.1.1"),
		Port: 8000,
	}
	seqnum1         uint64 = 12
	seqnum2         uint64 = 22
	tunnelEndpoint1        = "2.2.2.2:8000"
	tunnelPubkey1          = "derp"
	tunnelEndpoint2        = "3.3.3.3:8000"
	tunnelPubkey2          = "flerp"
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
			Seqnum:          acct.Seqnum,
			SourcePublicKey: acct.PublicKey,
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
	if msg.SourcePublicKey != *pubkey1 {
		t.Fatal("msg.PublicKey incorrect")
	}
	if msg.ControlAddress.String() != controlAddress1.String() {
		t.Fatal("msg.ControlAddress incorrect", msg.ControlAddress.String(), controlAddress1.String())
	}
	if msg.Seqnum != seqnum1 {
		t.Fatal("msg.Seqnum incorrect")
	}
	if msg.Confirm != confirm {
		t.Fatal("Confirm incorrect")
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
		Seqnum:    seqnum2,
		PublicKey: *pubkey2,
	}

	neighbor.Tunnel.Endpoint = tunnelEndpoint2
	neighbor.Tunnel.PublicKey = tunnelPubkey2

	msg := types.TunnelMessage{
		MessageMetadata: types.MessageMetadata{
			SourcePublicKey:      acct.PublicKey,
			DestinationPublicKey: neighbor.PublicKey,
			Seqnum:               acct.Seqnum,
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
	if msg.SourcePublicKey != *pubkey1 {
		t.Fatal("msg.PublicKey incorrect")
	}
	if msg.TunnelEndpoint != tunnelEndpoint2 {
		t.Fatal("msg.TunnelEndpoint incorrect", msg.TunnelEndpoint)
	}
	if msg.TunnelPublicKey != tunnelPubkey2 {
		t.Fatal("msg.TunnelPublicKey incorrect", msg.TunnelPublicKey)
	}
	if msg.Seqnum != seqnum1 {
		t.Fatal("msg.Seqnum incorrect")
	}
	if msg.Confirm != confirm {
		t.Fatal("Confirm incorrect")
	}

	var sig [ed25519.SignatureSize]byte

	if confirm {
		sig = [ed25519.SignatureSize]byte{0x9c, 0xa3, 0x94, 0xef, 0xf7, 0xd0, 0x2c, 0x58, 0x16, 0x7d, 0xfa, 0x7a, 0x55, 0xd, 0x7, 0xd9, 0x62, 0xb3, 0xeb, 0xc2, 0xfc, 0x34, 0x55, 0xf2, 0xec, 0x25, 0xbd, 0x92, 0xef, 0x33, 0xcb, 0xb0, 0x49, 0xa8, 0x6a, 0x66, 0x4e, 0x78, 0xda, 0x78, 0x98, 0xa8, 0x98, 0x19, 0xf, 0xa7, 0xa, 0x3f, 0xea, 0x79, 0x8f, 0x20, 0xbe, 0xbd, 0xb0, 0x2, 0x16, 0xf2, 0x23, 0xaf, 0xea, 0xfa, 0xaf, 0xf}
	} else {
		sig = [ed25519.SignatureSize]byte{0x91, 0xa9, 0x25, 0x9d, 0x73, 0x11, 0x1a, 0xab, 0x3f, 0xf4, 0x60, 0xa1, 0xf0, 0x3d, 0x13, 0x5a, 0x56, 0x11, 0x8c, 0x83, 0x17, 0x35, 0x2b, 0xa0, 0x25, 0x6, 0x2e, 0x1e, 0x89, 0x5f, 0x6b, 0x41, 0x11, 0xb0, 0x76, 0xcb, 0xa6, 0xcb, 0x42, 0x52, 0xce, 0xa4, 0xcd, 0x72, 0x12, 0xc1, 0xbe, 0x93, 0x5f, 0x61, 0x3f, 0xba, 0x5b, 0x56, 0x1e, 0x39, 0x3, 0x77, 0x5b, 0xe4, 0x2b, 0xa0, 0x6b, 0x4}
	}

	if msg.Signature != sig {
		t.Fatalf("msg.Signature incorrect: %#v SHOULD BE %#v", msg.Signature, sig)
	}
}
