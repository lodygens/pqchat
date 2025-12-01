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
package session

import "pqchat/src/internal/pqc"

// Session represents a symmetric session derived from a ML-KEM shared secret.
type Session struct {
	Cipher *pqc.AESGCM
}

// This encrypts an application message.
func (s *Session) Encrypt(plaintext []byte) ([]byte, error) {
	return s.Cipher.Encrypt(plaintext)
}

// This decrypts an application message.
func (s *Session) Decrypt(ciphertext []byte) ([]byte, error) {
	return s.Cipher.Decrypt(ciphertext)
}

// This returns a copy of the raw AES key (for debugging purposes).
func (s *Session) CipherKeyDebug() []byte {
	return s.Cipher.RawKey()
}
