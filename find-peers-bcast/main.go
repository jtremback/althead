package findPeersBCast

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/mdns"
)

type TOS struct {
	Speedtest string
	Result    string
	Denom     string
	Rate      int
}

func Advertise(
	tunnelIP net.IP,
	tunnelPort int,
	tunnelPubkey string,
	tos TOS,
) (*mdns.Server, error) {
	// Setup our service export
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(tos)
	if err != nil {
		return nil, err
	}

	service, err := mdns.NewMDNSService(
		tunnelPubkey,
		"Althea peer discovery",
		"local",
		host,
		tunnelPort,
		[]net.IP{tunnelIP},
		[]string{string(json)},
	)
	if err != nil {
		return nil, nil
	}

	// Create the mDNS server
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}

	return server, nil
}

func GetPeers() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	mdns.Lookup("Althea peer discovery", entriesCh)
	close(entriesCh)
}
