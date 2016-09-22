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

	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		fmt.Println(err)
	}

	if *server {
		fmt.Println("advertise")
		err := findPeersMCast.Advertise(
			iface,
			8481,
		)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("query")
		err := findPeersMCast.QueryPeers(
			iface,
			4500,
			8481,
		)
		if err != nil {
			fmt.Println(err)
		}
	}
}
