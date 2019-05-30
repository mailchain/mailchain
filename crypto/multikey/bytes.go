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

package multikey

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/pkg/errors"
)

type keyFunc func(data []byte) (crypto.PrivateKey, error)

// PrivateKeyFromBytes use the correct function to get the private key from bytes
func PrivateKeyFromBytes(keyType string, data []byte) (crypto.PrivateKey, error) {
	table := map[string]keyFunc{
		encoding.SECP256K1: func(data []byte) (crypto.PrivateKey, error) {
			return secp256k1.PrivateKeyFromBytes(data)
		},
	}

	keyFunc, ok := table[keyType]
	if !ok {
		return nil, errors.Errorf("unsupported curve type")
	}
	return keyFunc(data)
}
