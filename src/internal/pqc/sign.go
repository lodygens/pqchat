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

	pub, err = sig.GenerateKeyPair()
	if err != nil {
		return nil, nil, ErrSigGen
	}

	priv = sig.ExportSecretKey()
	return pub, priv, nil
}

// Sign returns a signature over the given message using the private key.
func Sign(message []byte, priv []byte) (sigBytes []byte, err error) {
	sig := oqs.Signature{}

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
