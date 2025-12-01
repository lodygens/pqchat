package pqc

import (
	"errors"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

var (
	ErrSigInit = errors.New("pqc: ML-DSA init failed")
	ErrSigGen  = errors.New("pqc: ML-DSA keygen failed")
	ErrSign    = errors.New("pqc: ML-DSA sign failed")
	ErrVerify  = errors.New("pqc: ML-DSA verify failed")
)

// ML-DSA-65 (Dilithium-3)
const SigAlgorithm = "ML-DSA-65"

// Generate signing keypair (public-key identity + private signing key).
func SigKeygen() (pub []byte, priv []byte, err error) {
	sig := oqs.Signature{}
	if err := sig.Init(SigAlgorithm, nil); err != nil {
		return nil, nil, ErrSigInit
	}
	defer sig.Clean()

	// GenerateKeyPair returns (publicKey, error), secret key is stored in the struct
	pub, err = sig.GenerateKeyPair()
	if err != nil {
		return nil, nil, ErrSigGen
	}

	// Export the secret key from the struct
	priv = sig.ExportSecretKey()
	return pub, priv, nil
}

// Sign returns a signature over the given message using the private key.
func Sign(message []byte, priv []byte) (sigBytes []byte, err error) {
	sig := oqs.Signature{}
	// Init takes (algorithm, secretKey) - pass the private key here
	if err := sig.Init(SigAlgorithm, priv); err != nil {
		return nil, ErrSign
	}
	defer sig.Clean()

	// Sign takes only the message, not the private key
	sigBytes, err = sig.Sign(message)
	if err != nil {
		return nil, ErrSign
	}
	return
}

// Verify checks if a signature is valid for (message, publicKey).
func Verify(message, signature, pub []byte) (bool, error) {
	sig := oqs.Signature{}
	if err := sig.Init(SigAlgorithm, nil); err != nil {
		return false, ErrSigInit
	}
	defer sig.Clean()

	ok, err := sig.Verify(message, signature, pub)
	if err != nil {
		return false, ErrVerify
	}
	return ok, nil
}
