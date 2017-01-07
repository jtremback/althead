package types

type Peer struct {
	ControlAddress string
	ControlPubkey  [32]byte
	TunnelAddress  string
	TunnelPubkey   []byte
}

type Account struct {
	ControlAddress string
	ControlPubkey  [32]byte
	ControlPrivkey [64]byte
	TunnelAddress  string
	TunnelPubkey   []byte
}
