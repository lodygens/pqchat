package session

import (
	"bufio"
	"fmt"

	p2pnet "github.com/libp2p/go-libp2p/core/network"

	"pqchat/src/internal/net"
	"pqchat/src/internal/pqc"
)

// This executes the ML-KEM handshake on the server side
// (the peer who receives the stream first).
func ServerHandshake(s p2pnet.Stream) (*Session, error) {
	// Generate the ML-KEM keypair
	kem, err := pqc.NewKEM()
	if err != nil {
		return nil, fmt.Errorf("new kem: %w", err)
	}
	defer kem.Clean()

	// Keygen returns (pub, priv, err) - priv is stored in kem and used by Decapsulate
	pub, _, err := kem.Keygen()
	if err != nil {
		return nil, fmt.Errorf("kem keygen: %w", err)
	}

	// Send the public key to the client
	if err := net.WriteFrame(s, pub); err != nil {
		return nil, fmt.Errorf("send pub: %w", err)
	}

	// Receive the ciphertext from the client
	rd := bufio.NewReader(s)
	ct, err := net.ReadFrame(rd)
	if err != nil {
		return nil, fmt.Errorf("recv ct: %w", err)
	}

	// Decapsulate the shared secret (uses the priv key stored in kem)
	ss, err := kem.Decapsulate(ct)
	if err != nil {
		return nil, fmt.Errorf("kem decaps: %w", err)
	}

	// Derive the AES-GCM key
	key, err := pqc.DeriveKey(ss, []byte("pqchat-handshake"))
	if err != nil {
		return nil, fmt.Errorf("derive key: %w", err)
	}

	// Create the AES-GCM instance
	aesgcm, err := pqc.NewAESGCM(key)
	if err != nil {
		return nil, fmt.Errorf("new aesgcm: %w", err)
	}

	return &Session{Cipher: aesgcm}, nil
}

// ClientHandshake executes the ML-KEM handshake on the client side
// (the peer who initiates the stream).
func ClientHandshake(s p2pnet.Stream) (*Session, error) {
	rd := bufio.NewReader(s)

	// 1. Receive the server's public key
	pub, err := net.ReadFrame(rd)
	if err != nil {
		return nil, fmt.Errorf("recv pub: %w", err)
	}

	// 2. Encapsulate → ciphertext + shared secret
	ct, ss, err := pqc.Encapsulate(pub)
	if err != nil {
		return nil, fmt.Errorf("kem encaps: %w", err)
	}

	// 3. Send the ciphertext to the server
	if err := net.WriteFrame(s, ct); err != nil {
		return nil, fmt.Errorf("send ct: %w", err)
	}

	// 4. Dérive la clé AES-GCM à partir du shared secret
	key, err := pqc.DeriveKey(ss, []byte("pqchat-handshake"))
	if err != nil {
		return nil, fmt.Errorf("derive key: %w", err)
	}

	aesgcm, err := pqc.NewAESGCM(key)
	if err != nil {
		return nil, fmt.Errorf("new aesgcm: %w", err)
	}

	return &Session{Cipher: aesgcm}, nil
}
