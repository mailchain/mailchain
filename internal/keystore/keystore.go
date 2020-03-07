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

package keystore

//go:generate mockgen -source=keystore.go -package=keystoretest -destination=./keystoretest/keystore_mock.go
import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
)

// Store private keys but does not return them, instead return decrypter or signer.
type Store interface {
	GetSigner(address []byte, protocol, network string, kind string, deriveKeyOptions multi.OptionsBuilders) (signer.Signer, error)
	GetDecrypter(address []byte, protocol, network string, decrypterType byte, deriveKeyOptions multi.OptionsBuilders) (cipher.Decrypter, error)
	Store(private crypto.PrivateKey, deriveKeyOptions multi.OptionsBuilders) (crypto.PublicKey, error)
	HasAddress(searchAddress []byte, protocol, network string) bool
	GetAddresses(protocol, address string) ([][]byte, error)
	GetPublicKeys() ([]crypto.PublicKey, error)
}
