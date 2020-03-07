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

package nacl

import (
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/pkg/errors"
)

// GetSigner return a transaction signer based on the supplied address.
func (f FileStore) GetSigner(address []byte, protocol, network string, kind string, deriveKeyOptions multi.OptionsBuilders) (signer.Signer, error) {
	encryptedKey, err := f.getEncryptedKeyByAddress(address, protocol, network)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pk, err := f.getPrivateKey(encryptedKey, deriveKeyOptions)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return keystore.Signer(protocol, kind, pk)
}
