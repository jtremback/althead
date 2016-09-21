package main

import (
	"fmt"
	"net"

	flags "github.com/jessevdk/go-flags"
	"github.com/jtremback/althea/find-peers-mdns"
)

var opts struct {
	server bool `short:"s" long:"server" description:"Run uplink advertisement server"`
}

// "github.com/jtremback/althea/find-peers-babel"

// func main() {
// 	fmt.Println("hello")
// 	_, err := findPeersBabel.Find(8481)
// 	fmt.Println(err)
// }

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
	}

	if opts.server {
		service := &findPeersMDNS.Service{
			Denom:        "ETH",
			Rate:         1,
			TunnelIP:     net.ParseIP("2000::1"),
			TunnelPort:   3456,
			TunnelPubkey: "shibb",
		}
		server, err := findPeersMDNS.Advertise(service)

		if err != nil {
			fmt.Println(err)
		}
		defer server.Shutdown()
	}
	findPeersMDNS.GetPeers()
}
