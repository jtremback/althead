package wireguard

import (
	"fmt"
	"testing"
)

func TestParseConfig(t *testing.T) {
	testConfig := `// [Interface]
// PrivateKey = yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=
// ListenPort = 51820

// [Peer]
// PublicKey = xTIBA5rboUvnH4htodjb6e697QjLERt1NAB4mZqp8Dg=
// Endpoint = 192.95.5.67:1234
// AllowedIPs = 10.192.122.3/32, 10.192.124.1/24`

	config, err := ParseConfig(testConfig)
	if err != nil {
		t.Fatal(err)
	}

	if config.PrivateKey != "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=" {
		t.Fatal("config.PrivateKey incorrect")
	}

	if config.ListenPort != 51820 {
		t.Fatal("config.ListenPort incorrect")
	}

	if config.Peer.AllowedIPs != "10.192.122.3/32, 10.192.124.1/24" {
		t.Fatal("config.Peer.AllowedIPs[0] incorrect: ", config.Peer.AllowedIPs)
	}

	fmt.Println(config.Peer.AllowedIPs[0])
}
