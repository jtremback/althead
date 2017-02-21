package serialization

import (
	"testing"

	"strings"

	"github.com/agl/ed25519"
	"github.com/incentivized-mesh-infrastructure/scrooge/types"
)

var (
	pubkey1                     = &[ed25519.PublicKeySize]byte{44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	privkey1                    = &[ed25519.PrivateKeySize]byte{112, 69, 149, 144, 72, 233, 25, 188, 124, 215, 67, 200, 213, 237, 133, 127, 215, 253, 230, 134, 26, 202, 25, 214, 36, 19, 233, 87, 212, 169, 119, 226, 44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	pubkey2                     = &[ed25519.PublicKeySize]byte{175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123, 162}
	privkey2                    = &[ed25519.PrivateKeySize]byte{13, 170, 251, 93, 50, 201, 207, 72, 224, 172, 35, 48, 16, 245, 116, 20, 88, 33, 155, 12, 226, 126, 59, 36, 184, 111, 95, 87, 156, 104, 140, 243, 175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123}
	helloMessage                = "scrooge_hello LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= 12 vi+JhQG9/u/aoRE5XQSQEpj6vyHoB+FaQQ2bwSHRBBgsb3KrmH71rrbcjbCOMYjJjnfvBVDAhqonodMo0sI0DQ=="
	helloConfirmMessage         = "scrooge_hello_confirm LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= 12 rXeNgP6HXpyLfkeW00bmfkclqjWE/wv2TliIZkUqEcB7cwT46WoeCo/2JYejaOHz+x1IFn/eq7hK8y+QRol9Cw=="
	tunnelMessage               = "scrooge_tunnel LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= r24MX1Kp720po7ddTcUjKSPLXsjYBimBqgwIYdMce6I= flerp 3.3.3.3:8000 12 kaklnXMRGqs/9GCh8D0TWlYRjIMXNSugJQYuHolfa0ERsHbLpstCUs6kzXISwb6TX2E/ultWHjkDd1vkK6BrBA=="
	tunnelConfirmMessage        = "scrooge_tunnel_confirm LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= r24MX1Kp720po7ddTcUjKSPLXsjYBimBqgwIYdMce6I= flerp 3.3.3.3:8000 12 nKOU7/fQLFgWffp6VQ0H2WKz68L8NFXy7CW9ku8zy7BJqGpmTnjaeJiomBkPpwo/6nmPIL69sAIW8iOv6vqvDw=="
	iface1                      = "eth0"
	seqnum1              uint64 = 12
	seqnum2              uint64 = 22
	tunnelEndpoint1             = "2.2.2.2:8000"
	tunnelPubkey1               = "derp"
	tunnelEndpoint2             = "3.3.3.3:8000"
	tunnelPubkey2               = "flerp"
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
		Confirm: confirm,
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
	if msg.Seqnum != seqnum1 {
		t.Fatal("msg.Seqnum incorrect")
	}
	if msg.Confirm != confirm {
		t.Fatal("Confirm incorrect")
	}

	var sig [ed25519.SignatureSize]byte

	if confirm {
		sig = [ed25519.SignatureSize]byte{0xad, 0x77, 0x8d, 0x80, 0xfe, 0x87, 0x5e, 0x9c, 0x8b, 0x7e, 0x47, 0x96, 0xd3, 0x46, 0xe6, 0x7e, 0x47, 0x25, 0xaa, 0x35, 0x84, 0xff, 0xb, 0xf6, 0x4e, 0x58, 0x88, 0x66, 0x45, 0x2a, 0x11, 0xc0, 0x7b, 0x73, 0x4, 0xf8, 0xe9, 0x6a, 0x1e, 0xa, 0x8f, 0xf6, 0x25, 0x87, 0xa3, 0x68, 0xe1, 0xf3, 0xfb, 0x1d, 0x48, 0x16, 0x7f, 0xde, 0xab, 0xb8, 0x4a, 0xf3, 0x2f, 0x90, 0x46, 0x89, 0x7d, 0xb}
	} else {
		sig = [ed25519.SignatureSize]byte{0xbe, 0x2f, 0x89, 0x85, 0x1, 0xbd, 0xfe, 0xef, 0xda, 0xa1, 0x11, 0x39, 0x5d, 0x4, 0x90, 0x12, 0x98, 0xfa, 0xbf, 0x21, 0xe8, 0x7, 0xe1, 0x5a, 0x41, 0xd, 0x9b, 0xc1, 0x21, 0xd1, 0x4, 0x18, 0x2c, 0x6f, 0x72, 0xab, 0x98, 0x7e, 0xf5, 0xae, 0xb6, 0xdc, 0x8d, 0xb0, 0x8e, 0x31, 0x88, 0xc9, 0x8e, 0x77, 0xef, 0x5, 0x50, 0xc0, 0x86, 0xaa, 0x27, 0xa1, 0xd3, 0x28, 0xd2, 0xc2, 0x34, 0xd}
	}

	if msg.Signature != sig {
		t.Fatalf("msg.Signature incorrect: %#v SHOULD BE %#v", msg.Signature, sig)
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
