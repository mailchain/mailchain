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

package aes256cbc

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
)

type invalidPrivateKey struct {
}

func (i invalidPrivateKey) Bytes() []byte             { return []byte{} }
func (i invalidPrivateKey) PublicKey() crypto.PublicKey { return nil }

type invalidPublicKey struct {
}

func (i invalidPublicKey) Bytes() []byte   { return []byte{} }
func (i invalidPublicKey) Address() []byte { return nil }
func Test_asPrivateECIES(t *testing.T) {
	type args struct {
		pk crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success-secp256k1-val",
			args{
				func() secp256k1.PrivateKey {
					k, err := secp256k1.PrivateKeyFromHex("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
					if err != nil {
						t.Error(err)
					}
					return *k
				}(),
			},
			false,
			false,
		},
		{
			"success-secp256k1-pointer",
			args{
				func() *secp256k1.PrivateKey {
					k, err := secp256k1.PrivateKeyFromHex("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
					if err != nil {
						t.Error(err)
					}
					return k
				}(),
			},
			false,
			false,
		},
		{
			"success-secp256k1-pointer",
			args{
				func() *secp256k1.PrivateKey {
					k, err := secp256k1.PrivateKeyFromHex("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
					if err != nil {
						t.Error(err)
					}
					return k
				}(),
			},
			false,
			false,
		},
		{
			"err-unsupported",
			args{
				func() invalidPrivateKey {
					return invalidPrivateKey{}
				}(),
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := asPrivateECIES(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("asPrivateECIES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantNil {
				t.Errorf("asPrivateECIES() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_asPublicECIES(t *testing.T) {
	type args struct {
		pk crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success-val",
			args{
				func() crypto.PublicKey {
					k, err := secp256k1.PrivateKeyFromHex("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
					if err != nil {
						t.Error(err)
					}
					return k.PublicKey()
				}(),
			},
			false,
			false,
		},
		{
			"err-invalid",
			args{
				func() invalidPublicKey {
					return invalidPublicKey{}
				}(),
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := asPublicECIES(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("asPublicECIES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("asPublicECIES() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}
