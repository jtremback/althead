package neighborAPI

import (
	"net"
	"testing"

	"fmt"

	"github.com/agl/ed25519"
	"github.com/incentivized-mesh-infrastructure/scrooge/types"
)

var (
	pubkey1                     = &[ed25519.PublicKeySize]byte{44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	privkey1                    = &[ed25519.PrivateKeySize]byte{112, 69, 149, 144, 72, 233, 25, 188, 124, 215, 67, 200, 213, 237, 133, 127, 215, 253, 230, 134, 26, 202, 25, 214, 36, 19, 233, 87, 212, 169, 119, 226, 44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	pubkey2                     = &[ed25519.PublicKeySize]byte{175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123, 162}
	privkey2                    = &[ed25519.PrivateKeySize]byte{13, 170, 251, 93, 50, 201, 207, 72, 224, 172, 35, 48, 16, 245, 116, 20, 88, 33, 155, 12, 226, 126, 59, 36, 184, 111, 95, 87, 156, 104, 140, 243, 175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123}
	helloMessage                = "scrooge_hello LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= 1.1.1.1:8000 12 bLKZBGYmwOiKnFflOEuZ0iRmIA0YtPeAyJM/wyh10pPeH+gYG2tMdDDCdhwYQKJtB3CTZ+Mn9s5IfxnbQB7jCQ=="
	helloConfirmMessage         = "scrooge_hello_confirm LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= 1.1.1.1:8000 12 WV4qgi9H5s5mbXg5B1LvKHgXepm+HkZmymqavQza7+1CvTLb+6R3yAJ1GMFPCetPfX8GLNTS4AJiN2tBhf7WAg=="
	tunnelMessage               = "scrooge_tunnel LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= derp 2.2.2.2:8000 12 6ylkg9zJjk0yYvudq2mUQVi+CaFo9faLsT8Yjiy/g95mTzbrul3eUJabrdiRm/gnbyS5lQHff5TjI2DMfaOZBw=="
	tunnelConfirmMessage        = "scrooge_tunnel_confirm LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= derp 2.2.2.2:8000 12 DEXqnvdy8ppaulkXO7ivJr+bQphGiSL3zAGFuMgvg4wvbuUlcAX+3+6bYE2xwssUmjtCyK7Z+XkKGhzHC5ChBw=="
	iface1                      = "eth0"
	controlAddress1             = "1.1.1.1:8000"
	controlAddress2             = "3.3.3.3:8000"
	seqnum1              uint64 = 12
	tunnelEndpoint1             = "2.2.2.2:8000"
	tunnelPubkey1               = "derp"
)

type fakeNetwork struct {
	SendUDPArgs interface{}
}

func (self *fakeNetwork) SendUDP(addr *net.UDPAddr, s string) error {
	self.SendUDPArgs = struct {
		addr net.UDPAddr
		s    string
	}{
		*addr,
		s,
	}
	return nil
}

func (self *fakeNetwork) SendMulticastUDP(*net.Interface, int, string) error {
	return nil
}

func TestHelloHandler(t *testing.T) {

	fn := &fakeNetwork{}

	n := NeighborAPI{
		Neighbors: map[[ed25519.PublicKeySize]byte]*types.Neighbor{},
		Account: &types.Account{
			PublicKey:        *pubkey1,
			PrivateKey:       *privkey1,
			ControlAddresses: map[string]string{},
			Seqnum:           12,
		},
		Network: fn,
	}

	n.Neighbors[*pubkey2] = &types.Neighbor{}
	n.Account.ControlAddresses[iface1] = controlAddress1

	err := n.Handlers([]byte(helloMessage), iface1)
	if err != nil {
		t.Error(err)
	}

	sendUDPArgs := fmt.Sprint(fn.SendUDPArgs)
	correctSendUDPArgs := "{{[0 0 0 0 0 0 0 0 0 0 255 255 1 1 1 1] 8000 } " + helloConfirmMessage + "}"

	if sendUDPArgs != correctSendUDPArgs {
		t.Error("fn.SendUDPArgs incorrect: ", sendUDPArgs+" SHOULD BE "+correctSendUDPArgs)
	}
}
