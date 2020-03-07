package anysender

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/address/addresstest"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/stretchr/testify/assert"
)

func TestSigner_Sign(t *testing.T) {
	type fields struct {
		privateKey crypto.PrivateKey
	}
	type args struct {
		opts signer.Options
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantSignedTransaction interface{}
		wantErr               bool
	}{
		{
			"from-charlotte-ab-ending",
			fields{
				secp256k1test.CharlottePrivateKey,
			},
			args{
				opts: SignerOptions{
					to:            addresstest.EthereumSofia,
					from:          addresstest.EthereumCharlotte,
					data:          encodingtest.MustDecodeHexZeroX("0x6d61696c636861696e010a82012e42d611d3c068ba7809d1f987a6b2881203d58cb2ec0401316e7633738d9110fec570cf76372908be72fb75a50e0754c16fac8a10670d3d7c75c1e210b607020ecdb4b5ce7a7360c1004c79a5d520a54c0592fb8acd2fdca5d735dd53d3ffef837d0f0f2623ca2cebdd587838acf5a0144541336db6ed6b4ab0981b893b91df50a0"),
					deadline:      7426351,
					refund:        500000000,
					gas:           100000,
					relayContract: encodingtest.MustDecodeHexZeroX("0xe8468689AB8607fF36663EE6522A7A595Ed8bC0C"),
				},
			},
			[]byte{0x14, 0xea, 0xc, 0xd0, 0xb4, 0x7, 0xe1, 0x4c, 0xd2, 0x1f, 0x30, 0x7d, 0xf0, 0xed, 0x9f, 0x6d, 0x70, 0xb4, 0xdb, 0x2c, 0xc4, 0x4d, 0xd, 0x1d, 0xec, 0xee, 0xc2, 0xa2, 0xfa, 0x94, 0x7f, 0xbb, 0x20, 0x16, 0x80, 0x48, 0xa1, 0xcd, 0xcc, 0x33, 0x3f, 0x19, 0x96, 0x30, 0x8e, 0xe5, 0xfb, 0xaf, 0x5b, 0xd9, 0xa6, 0x23, 0x60, 0x28, 0xd1, 0x24, 0xbd, 0x15, 0xb5, 0x9e, 0xf0, 0x6e, 0xc, 0x84, 0x1b},
			false,
		},
		{
			"from-charlotte-1c-ending",
			fields{
				secp256k1test.CharlottePrivateKey,
			},
			args{
				opts: SignerOptions{
					to:            addresstest.EthereumSofia,
					from:          addresstest.EthereumCharlotte,
					data:          encodingtest.MustDecodeHexZeroX("0x6d61696c636861696e010a82012e28443acc17cb3e6e5cfecf755d831abc02fb5e121c25cf35afae375541d019e6ff87af1345ec556aab929a950272faeeebbf8ec5228e06c940eb58f479f7a91f87e2e71e60d5211c60b55f7a73215824050f7477470e6cbef9f9307145595a4b0dd73d44b1ec9353d66b3f876c58c160897597758e978322a0e5d6190e3efea90c"),
					deadline:      7428265,
					refund:        500000000,
					gas:           83288,
					relayContract: encodingtest.MustDecodeHexZeroX("0xe8468689AB8607fF36663EE6522A7A595Ed8bC0C"),
				},
			},
			encodingtest.MustDecodeHexZeroX("0xe8cec4bcced88b310c85523e0d8a94d78d70f3261df62c322c652095671180542214a254eb165005add1502baab3aa5c87314cb3fc8ecd7325a020fbf5f2b2491c"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Signer{
				privateKey: tt.fields.privateKey,
			}
			gotSignedTransaction, err := e.Sign(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signer.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantSignedTransaction, gotSignedTransaction) {
				t.Errorf("Signer.Sign() = %v, want %v", gotSignedTransaction, tt.wantSignedTransaction)
			}
		})
	}
}
