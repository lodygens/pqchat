// Copyright 2025 Oleg Lodygensky
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions AND
// limitations under the License.

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	libhost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"

	p2pnet "pqchat/src/internal/net"
	"pqchat/src/internal/session"
)

var (
	flagRelay   = flag.String("relay", "", "relay multiaddr, e.g. /ip4/1.2.3.4/tcp/4001/p2p/<id>")
	flagConnect = flag.String("connect", "", "peer multiaddr to connect to (optional)")
)

func main() {
	flag.Parse()

	if *flagRelay == "" {
		fmt.Println("⚠️ No relay configured, running in direct TCP mode.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// gestion Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n[!] Interrupt, shutting down…")
		cancel()
	}()

	h, err := p2pnet.NewHost("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		fmt.Println("Cannot create host:", err)
		return
	}

	fmt.Println("Local PeerID:", h.ID())
	fmt.Println("Listening on:")
	for _, a := range h.Addrs() {
		fmt.Println("   ", a)
	}

	if err := connectRelay(ctx, h, *flagRelay); err != nil {
		fmt.Println("Relay connect failed:", err)
		return
	}

	// Handler for incoming streams
	h.SetStreamHandler("/pqchat/1.0.0", func(s network.Stream) {
		fmt.Println("\nIncoming connection from", s.Conn().RemotePeer())

		sess, err := session.ServerHandshake(s)
		if err != nil {
			fmt.Println("Handshake (server) failed:", err)
			_ = s.Reset()
			return
		}
		fmt.Println("PQC session established (server side)")

		rd := bufio.NewReader(s)

		for {
			frame, err := p2pnet.ReadFrame(rd)
			if err != nil {
				// stream fermé
				return
			}

			pt, err := sess.Decrypt(frame)
			if err != nil {
				fmt.Println("Decrypt failed:", err)
				return
			}

			fmt.Printf("\n[%s] %s\n> ", s.Conn().RemotePeer(), string(pt))
		}
	})

	// If we have a destination peer, connect to it immediately
	var (
		activeSess *session.Session
		activeStrm network.Stream
	)

	if *flagConnect != "" {
		activeSess, activeStrm, err = connectToPeer(ctx, h, *flagConnect)
		if err != nil {
			fmt.Println("Initial peer connect failed:", err)
		}
	}

	// Interactive loop
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			fmt.Print("> ")
			continue
		}

		// If no active session, try to connect now
		if activeSess == nil || activeStrm == nil {
			if *flagConnect == "" {
				fmt.Println("No peer connected. Start pqchat with -connect <multiaddr>")
				fmt.Print("> ")
				continue
			}
			activeSess, activeStrm, err = connectToPeer(ctx, h, *flagConnect)
			if err != nil {
				fmt.Println("Cannot connect to peer:", err)
				fmt.Print("> ")
				continue
			}
		}

		ct, err := activeSess.Encrypt([]byte(line))
		if err != nil {
			fmt.Println("Encrypt failed:", err)
			fmt.Print("> ")
			continue
		}

		if err := p2pnet.WriteFrame(activeStrm, ct); err != nil {
			fmt.Println("Send failed:", err)
			activeSess = nil
			activeStrm = nil
			fmt.Print("> ")
			continue
		}

		fmt.Print("> ")
	}

	// fin
	_ = h.Close()
}

/* -----------------------------------------------------------
   Connexion au relay
-----------------------------------------------------------*/

func connectRelay(ctx context.Context, h libhost.Host, maddrStr string) error {
	if maddrStr == "" {
		return nil // skip if no relay (for local test)
	}

	info, err := peer.AddrInfoFromString(maddrStr)
	if err != nil {
		return fmt.Errorf("invalid relay multiaddr: %w", err)
	}

	fmt.Println("Connecting to relay:", info.ID)
	if err := h.Connect(ctx, *info); err != nil {
		return fmt.Errorf("relay connect: %w", err)
	}
	fmt.Println("Relay connected")
	return nil
}

/* -----------------------------------------------------------
   Connexion à un peer + handshake PQC côté client
-----------------------------------------------------------*/

func connectToPeer(ctx context.Context, h libhost.Host, maddrStr string) (*session.Session, network.Stream, error) {
	info, err := peer.AddrInfoFromString(maddrStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid peer multiaddr: %w", err)
	}

	fmt.Println("Connecting to peer:", info.ID)
	if err := h.Connect(ctx, *info); err != nil {
		return nil, nil, fmt.Errorf("peer connect: %w", err)
	}

	s, err := h.NewStream(ctx, info.ID, "/pqchat/1.0.0")
	if err != nil {
		return nil, nil, fmt.Errorf("open stream: %w", err)
	}

	fmt.Println("Running PQC client handshake…")
	sess, err := session.ClientHandshake(s)
	if err != nil {
		_ = s.Reset()
		return nil, nil, fmt.Errorf("pqc handshake: %w", err)
	}

	fmt.Println("PQC session established (client side)")
	return sess, s, nil
}
