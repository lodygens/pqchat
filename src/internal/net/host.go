package net

import (
	libp2p "github.com/libp2p/go-libp2p"
	noisy "github.com/libp2p/go-libp2p/p2p/security/noise"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	tcp "github.com/libp2p/go-libp2p/p2p/transport/tcp"

	"github.com/libp2p/go-libp2p/core/host"
)

func NewHost(listen string) (host.Host, error) {
	return libp2p.New(
		libp2p.Security(noisy.ID, noisy.New),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(quic.NewTransport),
		libp2p.ListenAddrStrings(listen),
		libp2p.EnableRelay(),
	)
}
