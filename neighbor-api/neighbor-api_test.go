package neighborAPI

import (
	"net"
	"testing"

	"fmt"

	"github.com/agl/ed25519"
	"github.com/incentivized-mesh-infrastructure/scrooge/serialization"
	"github.com/incentivized-mesh-infrastructure/scrooge/types"
)

var (
	pubkey1         = &[ed25519.PublicKeySize]byte{44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	privkey1        = &[ed25519.PrivateKeySize]byte{112, 69, 149, 144, 72, 233, 25, 188, 124, 215, 67, 200, 213, 237, 133, 127, 215, 253, 230, 134, 26, 202, 25, 214, 36, 19, 233, 87, 212, 169, 119, 226, 44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	pubkey2         = &[ed25519.PublicKeySize]byte{175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123, 162}
	privkey2        = &[ed25519.PrivateKeySize]byte{13, 170, 251, 93, 50, 201, 207, 72, 224, 172, 35, 48, 16, 245, 116, 20, 88, 33, 155, 12, 226, 126, 59, 36, 184, 111, 95, 87, 156, 104, 140, 243, 175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123}
	iface           = "eth0"
	controlAddress1 = "1.1.1.1:8000"
	controlAddress2 = "2.2.2.2:8000"
	account1        = &types.Account{
		PublicKey:  [ed25519.PublicKeySize]byte{44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85},
		PrivateKey: [ed25519.PrivateKeySize]byte{112, 69, 149, 144, 72, 233, 25, 188, 124, 215, 67, 200, 213, 237, 133, 127, 215, 253, 230, 134, 26, 202, 25, 214, 36, 19, 233, 87, 212, 169, 119, 226, 44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85},
		ControlAddresses: map[string]string{
			(iface): "1.1.1.1:8000",
		},
		Seqnum: 16,
	}
	account2 = &types.Account{
		PublicKey:  [ed25519.PublicKeySize]byte{175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123, 162},
		PrivateKey: [ed25519.PrivateKeySize]byte{13, 170, 251, 93, 50, 201, 207, 72, 224, 172, 35, 48, 16, 245, 116, 20, 88, 33, 155, 12, 226, 126, 59, 36, 184, 111, 95, 87, 156, 104, 140, 243, 175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123},
		ControlAddresses: map[string]string{
			(iface): "2.2.2.2:8000",
		},
		Seqnum: 16,
	}
)

type fakeNetwork struct {
	SendUDPArgs string
}

func (self *fakeNetwork) SendUDP(addr *net.UDPAddr, s string) error {
	self.SendUDPArgs = fmt.Sprint(addr, s)
	return nil
}

func (self *fakeNetwork) SendMulticastUDP(*net.Interface, int, string) error {
	return nil
}

// TestHelloHandler simulates account2 receiving and responding to a
// scrooge_hello message from account1
func TestHelloHandler(t *testing.T) {

	fn := &fakeNetwork{}

	n := NeighborAPI{
		Neighbors: map[[ed25519.PublicKeySize]byte]*types.Neighbor{},
		Account:   account2,
		Network:   fn,
	}

	helloMessage, err := serialization.FmtHello(account1, controlAddress1, false)
	if err != nil {
		t.Fatal(err)
	}

	err = n.Handlers([]byte(helloMessage), iface)
	if err != nil {
		t.Fatal(err)
	}

	helloConfirmMessage, err := serialization.FmtHello(account2, controlAddress2, true)
	if err != nil {
		t.Fatal(err)
	}

	correctSendUDPArgs := controlAddress1 + helloConfirmMessage

	if fn.SendUDPArgs != correctSendUDPArgs {
		t.Fatal("fn.SendUDPArgs incorrect: ", fn.SendUDPArgs, " SHOULD BE ", correctSendUDPArgs)
	}

	if n.Account.Seqnum != 17 {
		t.Fatal("n.Account.Seqnum incorrect: ", n.Account.Seqnum, " SHOULD BE ", 17)
	}
}

func TestBadSeqnum(t *testing.T) {
	fn := &fakeNetwork{}

	n2 := NeighborAPI{
		Neighbors: map[[ed25519.PublicKeySize]byte]*types.Neighbor{},
		Account:   account2,
		Network:   fn,
	}

	// We receive a message
	helloMessage, err := serialization.FmtHello(account1, controlAddress1, false)
	if err != nil {
		t.Fatal(err)
	}

	err = n2.Handlers([]byte(helloMessage), iface)
	if err != nil {
		t.Fatal(err)
	}

	// Now we receive the same message without incrementing the sequence number
	helloMessage, err = serialization.FmtHello(account1, controlAddress1, false)
	if err != nil {
		t.Fatal(err)
	}

	err = n2.Handlers([]byte(helloMessage), iface)
	if err != nil {
		if err.Error() != "sequence number too low" {
			t.Fatal("wrong error for bad sequence number ", err)
		}
	} else {
		t.Fatal("no error for bad sequence number")
	}

	// Now we do increment the sequence number
	account1.Seqnum = account1.Seqnum + 1
	helloMessage, err = serialization.FmtHello(account1, controlAddress1, false)
	if err != nil {
		t.Fatal(err)
	}

	err = n2.Handlers([]byte(helloMessage), iface)
	if err != nil {
		t.Fatal(err)
	}
}
