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
package pqc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/hkdf"
)

const (
	AESKeySize   = 32 // 256-bit
	AESNonceSize = 12 // GCM standard nonce
)

// AESGCM wraps a cipher.AEAD for easier use.
type AESGCM struct {
	aead cipher.AEAD
	key  []byte // Store the raw key for debug access
}

// DeriveKey derives a 256-bit AES key from a ML-KEM shared secret
// using HKDF-SHA256. contextInfo can be "pqchat-session", peer IDs, etc.
func DeriveKey(sharedSecret []byte, contextInfo []byte) ([]byte, error) {
	if len(sharedSecret) == 0 {
		return nil, errors.New("pqc: empty shared secret")
	}
	h := hkdf.New(sha256.New, sharedSecret, nil, contextInfo)
	key := make([]byte, AESKeySize)
	if _, err := io.ReadFull(h, key); err != nil {
		return nil, err
	}
	return key, nil
}

// NewAESGCM creates an AES-GCM instance from a raw key.
func NewAESGCM(key []byte) (*AESGCM, error) {
	if len(key) != AESKeySize {
		return nil, errors.New("pqc: invalid AES key size")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// Store a copy of the key for debug access
	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)
	return &AESGCM{aead: aead, key: keyCopy}, nil
}

// Encrypt encrypts plaintext and returns nonce || ciphertext.
func (a *AESGCM) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, AESNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	ct := a.aead.Seal(nil, nonce, plaintext, nil)
	out := make([]byte, 0, len(nonce)+len(ct))
	out = append(out, nonce...)
	out = append(out, ct...)
	return out, nil
}

// Decrypt expects nonce || ciphertext.
func (a *AESGCM) Decrypt(data []byte) ([]byte, error) {
	if len(data) < AESNonceSize {
		return nil, errors.New("pqc: ciphertext too short")
	}
	nonce := data[:AESNonceSize]
	ct := data[AESNonceSize:]
	return a.aead.Open(nil, nonce, ct, nil)
}

// RawKey returns a copy of the raw AES key (for debugging purposes).
func (a *AESGCM) RawKey() []byte {
	if a.key == nil {
		return nil
	}
	keyCopy := make([]byte, len(a.key))
	copy(keyCopy, a.key)
	return keyCopy
}
