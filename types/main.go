package types

import (
	"net"
)

type Neighbor struct {
	ControlAddress string
	ControlPubkey  []byte
	TunnelAddress  string
	TunnelPubkey   []byte
	Interface      *net.Interface
}

type Tunnel struct {
	ControlPubkey  []byte
	ControlPrivkey []byte
	ControlAddress string
	TunnelAddress  string
	TunnelPubkey   []byte
	TunnelPrivkey  []byte
	Interface      *net.Interface
}
