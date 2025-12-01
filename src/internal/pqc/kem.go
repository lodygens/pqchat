package pqc

import (
	"github.com/open-quantum-safe/liboqs-go/oqs"
)

const DefaultKEM = "ML-KEM-768"

// One object holds one complete ML-KEM instance
type KEM struct {
	obj *oqs.KeyEncapsulation
}

func NewKEM() (*KEM, error) {
	kem := &oqs.KeyEncapsulation{}
	if err := kem.Init(DefaultKEM, nil); err != nil {
		return nil, err
	}
	return &KEM{obj: kem}, nil
}

// Generate and keep the SAME KEM object for keygen
func (k *KEM) Keygen() (pub, priv []byte, err error) {
	pub, err = k.obj.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	priv = k.obj.ExportSecretKey()
	return pub, priv, nil
}

// Encapsulate: stateless, ok to use new object
func Encapsulate(peerPub []byte) (ct, ss []byte, err error) {
	k, err := NewKEM()
	if err != nil {
		return nil, nil, err
	}
	defer k.obj.Clean()

	// Try EncapSecret with the public key as parameter
	return k.obj.EncapSecret(peerPub)
}

// Decapsulate MUST use the SAME k.obj that created the priv
func (k *KEM) Decapsulate(ct []byte) ([]byte, error) {
	return k.obj.DecapSecret(ct)
}

// Clean cleans up the KEM object
func (k *KEM) Clean() {
	if k.obj != nil {
		k.obj.Clean()
	}
}
