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

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/sender/anysender"
	"github.com/pkg/errors"
)

// Signer use the correct function to get the decrypter from private key
func Signer(protocol string, kind string, pk crypto.PrivateKey) (signer.Signer, error) {
	switch protocol {
	case protocols.Ethereum:
		if kind == "any.sender" {
			return anysender.NewSigner(pk)
		}
		return ethereum.NewSigner(pk)
	default:
		return nil, errors.Errorf("unsupported signer type")
	}
}
