package types

type Peer struct {
	ControlAddress string
	ControlPubkey  []byte
	TunnelAddress  string
	TunnelPubkey   []byte
}

type Account struct {
	ControlAddress string
	ControlPubkey  []byte
	ControlPrivkey []byte
	TunnelAddress  string
	TunnelPubkey   []byte
}
