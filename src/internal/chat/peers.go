// Copyright 2025 Oleg Lodygensky
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions AND
// limitations under the License.
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
