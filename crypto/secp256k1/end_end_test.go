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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignVerify(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name         string
		signedBy     PrivateKey
		verifiedBy   PublicKey
		message      []byte
		wantSig      []byte
		wantErr      bool
		wantVerified bool
	}{
		{
			"charlotte-private-public-key",
			*charlottePrivateKey().(*PrivateKey),
			charlottePublicKey().(PublicKey),
			[]byte("message"),
			[]byte{0x9d, 0xf7, 0x76, 0xab, 0xde, 0x8c, 0x20, 0x55, 0xc3, 0x4, 0x68, 0x37, 0xa8, 0x66, 0xf8, 0x89, 0x95, 0xf9, 0x82, 0xf0, 0x4b, 0xb8, 0x23, 0x40, 0xf0, 0x3, 0x8, 0x6a, 0x32, 0xa7, 0xac, 0xef, 0x5f, 0xa, 0xea, 0xda, 0x60, 0xbf, 0x9, 0xd5, 0xc3, 0x27, 0x61, 0xa, 0xc5, 0xc8, 0x33, 0xe3, 0xa0, 0x79, 0xdf, 0x6d, 0xe1, 0x9c, 0xa8, 0xcc, 0x33, 0xea, 0x1d, 0xe6, 0x3, 0x34, 0xb1, 0xa1, 0x0},
			false,
			true,
		},
		{
			"sofia-private-public-key",
			*sofiaPrivateKey().(*PrivateKey),
			sofiaPublicKey().(PublicKey),
			[]byte("egassem"),
			[]byte{0xe9, 0x33, 0xe, 0x4a, 0xe3, 0x5, 0x19, 0xea, 0x36, 0x37, 0x19, 0xdd, 0xbc, 0x91, 0xfd, 0x4f, 0xd3, 0x64, 0x9b, 0xdc, 0xf0, 0x74, 0x36, 0x16, 0xc9, 0x81, 0xfc, 0x6d, 0x3c, 0x7e, 0xb0, 0xd0, 0x6e, 0xdd, 0x4, 0x13, 0xfd, 0x15, 0xe5, 0xec, 0x64, 0x6e, 0x63, 0xe0, 0x84, 0xdb, 0xb2, 0xd7, 0xcf, 0x18, 0x3d, 0x81, 0x1e, 0x31, 0x36, 0x77, 0x39, 0x86, 0x4b, 0x58, 0xb8, 0x23, 0xed, 0xc, 0x1},
			false,
			true,
		},
		{
			"sofia-private-charlotte-public-key",
			*sofiaPrivateKey().(*PrivateKey),
			charlottePublicKey().(PublicKey),
			[]byte("egassem"),
			[]byte{0xe9, 0x33, 0xe, 0x4a, 0xe3, 0x5, 0x19, 0xea, 0x36, 0x37, 0x19, 0xdd, 0xbc, 0x91, 0xfd, 0x4f, 0xd3, 0x64, 0x9b, 0xdc, 0xf0, 0x74, 0x36, 0x16, 0xc9, 0x81, 0xfc, 0x6d, 0x3c, 0x7e, 0xb0, 0xd0, 0x6e, 0xdd, 0x4, 0x13, 0xfd, 0x15, 0xe5, 0xec, 0x64, 0x6e, 0x63, 0xe0, 0x84, 0xdb, 0xb2, 0xd7, 0xcf, 0x18, 0x3d, 0x81, 0x1e, 0x31, 0x36, 0x77, 0x39, 0x86, 0x4b, 0x58, 0xb8, 0x23, 0xed, 0xc, 0x1},
			false,
			false,
		},
		{
			"charlotte-private-sofia-public-key",
			*charlottePrivateKey().(*PrivateKey),
			sofiaPublicKey().(PublicKey),
			[]byte("message"),
			[]byte{0x9d, 0xf7, 0x76, 0xab, 0xde, 0x8c, 0x20, 0x55, 0xc3, 0x4, 0x68, 0x37, 0xa8, 0x66, 0xf8, 0x89, 0x95, 0xf9, 0x82, 0xf0, 0x4b, 0xb8, 0x23, 0x40, 0xf0, 0x3, 0x8, 0x6a, 0x32, 0xa7, 0xac, 0xef, 0x5f, 0xa, 0xea, 0xda, 0x60, 0xbf, 0x9, 0xd5, 0xc3, 0x27, 0x61, 0xa, 0xc5, 0xc8, 0x33, 0xe3, 0xa0, 0x79, 0xdf, 0x6d, 0xe1, 0x9c, 0xa8, 0xcc, 0x33, 0xea, 0x1d, 0xe6, 0x3, 0x34, 0xb1, 0xa1, 0x0},
			false,
			false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotSig, err := tc.signedBy.Sign(tc.message)
			assert.Equal(tc.wantErr, err != nil)
			assert.Equal(tc.wantSig, gotSig)

			gotVerified := tc.verifiedBy.Verify(tc.message, gotSig)
			assert.Equal(tc.wantVerified, gotVerified)
		})
	}
}
