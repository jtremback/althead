package serialization

import (
	"crypto/rand"
	"testing"

	"fmt"

	"github.com/agl/ed25519"
	"github.com/jtremback/scrooge/types"
)

var (
	pubkey1      = &[ed25519.PublicKeySize]byte{44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	privkey1     = &[ed25519.PrivateKeySize]byte{112, 69, 149, 144, 72, 233, 25, 188, 124, 215, 67, 200, 213, 237, 133, 127, 215, 253, 230, 134, 26, 202, 25, 214, 36, 19, 233, 87, 212, 169, 119, 226, 44, 176, 80, 246, 247, 71, 5, 229, 108, 111, 158, 77, 18, 116, 98, 28, 84, 59, 215, 93, 182, 34, 240, 5, 147, 229, 211, 253, 44, 221, 237, 85}
	pubkey2      = &[ed25519.PublicKeySize]byte{175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123, 162}
	privkey2     = &[ed25519.PrivateKeySize]byte{13, 170, 251, 93, 50, 201, 207, 72, 224, 172, 35, 48, 16, 245, 116, 20, 88, 33, 155, 12, 226, 126, 59, 36, 184, 111, 95, 87, 156, 104, 140, 243, 175, 110, 12, 95, 82, 169, 239, 109, 41, 163, 183, 93, 77, 197, 35, 41, 35, 203, 94, 200, 216, 6, 41, 129, 170, 12, 8, 97, 211, 28, 123}
	helloMessage = "scrooge_hello LLBQ9vdHBeVsb55NEnRiHFQ71122IvAFk+XT/Szd7VU= derp.com 12 2RQBbXLrff2Bf0PMzkH6V33z/9NaaXxSAQ/IiHj9EntycInjpIYmIz2gJ7E/sf9PcKxr8OBffCiI31OrsHDiCQ=="
)

func TestFmtHello(t *testing.T) {
	pubkey, privkey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pubkey)
	fmt.Println(privkey)
	acct := &types.Account{
		Seqnum:     12,
		PublicKey:  *pubkey1,
		PrivateKey: *privkey1,
		ControlAddresses: map[string]string{
			"eth0": "derp.com",
		},
	}
	msg, err := FmtHello(acct, "eth0")
	if err != nil {
		t.Error(err)
	}

	if msg != helloMessage {
		t.Error("Message format incorrect: " + msg)
	}
	fmt.Println(msg)
}

func TestParseHello(t *testing.T) {

}
