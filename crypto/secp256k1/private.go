// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secp256k1

import (
	"crypto/ecdsa"
	"io"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/crypto"
	"github.com/pkg/errors"
)

// PrivateKey based on the secp256k1 curve
type PrivateKey struct {
	ecdsa ecdsa.PrivateKey
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	return ethcrypto.FromECDSA(&pk.ecdsa)
}

// Sign signs the message with the private key and returns the signature.
func (pk PrivateKey) Sign(hash []byte) (signature []byte, err error) {
	// hash := sha256.Sum256(message)
	return ethcrypto.Sign(hash, &pk.ecdsa)
}

// PublicKey return the public key that is derived from the private key
func (pk PrivateKey) PublicKey() crypto.PublicKey {
	return &PublicKey{ecdsa: pk.ecdsa.PublicKey}
}

// Kind is the type of private key.
func (pk PrivateKey) Kind() string {
	return crypto.KindSECP256K1
}

// ECIES returns an ECIES representation of the private key.
func (pk PrivateKey) ECIES() *ecies.PrivateKey {
	return ecies.ImportECDSA(&pk.ecdsa)
}

// ECDSA returns an ECDSA representation of the private key.
func (pk PrivateKey) ECDSA() (*ecdsa.PrivateKey, error) {
	rpk, err := ethcrypto.ToECDSA(pk.Bytes())
	return rpk, errors.WithMessage(err, "could not convert private key")
}

// PrivateKeyFromECDSA get a private key from an ecdsa.PrivateKey
func PrivateKeyFromECDSA(pk ecdsa.PrivateKey) PrivateKey {
	return PrivateKey{ecdsa: pk}
}

// PrivateKeyFromBytes get a private key from []byte
func PrivateKeyFromBytes(pk []byte) (*PrivateKey, error) {
	rpk, err := ethcrypto.ToECDSA(pk)
	if err != nil {
		return nil, errors.Errorf("could not convert private key")
	}

	return &PrivateKey{ecdsa: *rpk}, nil
}

func GenerateKey(rand io.Reader) (*PrivateKey, error) {
	pk, err := ecdsa.GenerateKey(ethcrypto.S256(), rand)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{ecdsa: *pk}, nil
}
