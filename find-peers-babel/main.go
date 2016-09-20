package findPeersBabel

import (
	"net"
	"net/rpc"
)

type TOS struct {
	QOSType      string
	QOSResult    string
	Denom        string
	Rate         int
	TunnelIP     net.IP
	TunnelPort   int
	TunnelPubkey string
}

type RPC struct {
	*TOS
}

func Advertise(
	port int,
	tos TOS,
) error {
	rpc.Register(&RPC{
		&TOS{
			QOSType:      "herp",
			QOSResult:    "derp",
			Denom:        "ETH",
			Rate:         0,
			TunnelIP:     net.ParseIP("::1"),
			TunnelPort:   8900,
			TunnelPubkey: "shib",
		},
	})

	l, err := net.Listen("udp", string(port))
	if err != nil {
		return err
	}

	rpc.Accept(l)
	return nil
}

func QueryTOS() error {

}
