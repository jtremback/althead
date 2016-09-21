package findPeersMDNS

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/mdns"
)

type Service struct {
	Denom        string
	Rate         int
	TunnelIP     net.IP
	TunnelPort   int
	TunnelPubkey string
}

func Advertise(
	service *Service,
) (*mdns.Server, error) {
	// Setup our service export
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(service)
	if err != nil {
		return nil, err
	}

	MDNSservice, err := mdns.NewMDNSService(
		service.TunnelPubkey,
		"Althea peer discovery",
		"local",
		host,
		service.TunnelPort,
		[]net.IP{service.TunnelIP},
		[]string{string(json)},
	)
	if err != nil {
		return nil, nil
	}

	// Create the mDNS server
	server, err := mdns.NewServer(&mdns.Config{Zone: MDNSservice})
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
