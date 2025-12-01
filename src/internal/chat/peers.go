package chat

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	userToPeer sync.Map // userID → peer.ID
)

func RegisterPeer(userID string, pid peer.ID) {
	userToPeer.Store(userID, pid)
}

func LookupPeer(userID string) (peer.ID, bool) {
	v, ok := userToPeer.Load(userID)
	if !ok {
		return "", false
	}
	return v.(peer.ID), true
}

func AllPeers() []peer.ID {
	var out []peer.ID
	userToPeer.Range(func(_, v any) bool {
		out = append(out, v.(peer.ID))
		return true
	})
	return out
}
