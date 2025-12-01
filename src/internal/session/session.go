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
