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

package address

import (
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/pkg/errors"
)

// EncodeByProtocol takes an address as `[]byte` then selects the relevant encoding method to encode it as string.
func EncodeByProtocol(in []byte, protocol string) (encoded, encodingType string, err error) {
	switch protocol {
	case protocols.Ethereum:
		encodingType = encoding.TypeHex0XPrefix
		encoded = encoding.EncodeHexZeroX(in)
	default:
		err = errors.Errorf("%q unsupported protocol", protocol)
	}

	return encoded, encodingType, err
}
