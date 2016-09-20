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

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
	}

	tos := findPeers.TOS{
		Speedtest: "",
		Result:    "",
		Denom:     "ETH",
		Rate:      0,
	}

	ip := net.ParseIP("10::1")
	if ip == nil {
		fmt.Println("fucked up ip address")
	}

	server, err := findPeersMDNS.Advertise(
		ip,
		7099,
		"shib",
		tos,
	)
	if err != nil {
		fmt.Println(err)
	}

	defer server.Shutdown()
}
