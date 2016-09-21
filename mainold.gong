package main

import

// flags "github.com/jessevdk/go-flags"

(
	"flag"
	"fmt"
	"net"

	"github.com/jtremback/althea/find-peers-mcast"
)

func main() {

	server := flag.Bool("s", false, "Run a server")

	flag.Parse()

	if *server {
		service := &findPeersMCast.Service{
			Denom:        "ETH",
			Rate:         1,
			TunnelIP:     net.ParseIP("2000::1"),
			TunnelPort:   3456,
			TunnelPubkey: "shibb",
		}
		err := findPeersMCast.Advertise(service)
		fmt.Println("called")
		fmt.Println(err)
	} else {
		err := findPeersMCast.GetPeers()
		fmt.Println(err)
	}
}

// "github.com/jtremback/althea/find-peers-babel"

// func main() {
// 	fmt.Println("hello")
// 	_, err := findPeersBabel.Find(8481)
// 	fmt.Println(err)
// }

// func main() {
// 	_, err := flags.Parse(&opts)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	if opts.server {
// 		service := &findPeersMDNS.Service{
// 			Denom:        "ETH",
// 			Rate:         1,
// 			TunnelIP:     net.ParseIP("2000::1"),
// 			TunnelPort:   3456,
// 			TunnelPubkey: "shibb",
// 		}
// 		server, err := findPeersMDNS.Advertise(service)

// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		defer server.Shutdown()
// 	}
// 	findPeersMDNS.GetPeers()
// }
