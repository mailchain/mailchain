package anysender

// func TestSender_encode(t *testing.T) {
// 	type fields struct {
// 		key            string
// 		networkConfigs map[string]networkConfig
// 	}
// 	type args struct {
// 		to            []byte
// 		from          []byte
// 		data          []byte
// 		deadline      int64
// 		refund        int64
// 		gas           int64
// 		relayContract []byte
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    []byte
// 		wantErr bool
// 	}{
// 		{
// 			"charlotte",
// 			fields{},
// 			args{
// 				addresstest.EthereumSofia,
// 				addresstest.EthereumCharlotte,
// 				encodingtest.MustDecodeHexZeroX("0x6d61696c636861696e010a82012e42d611d3c068ba7809d1f987a6b2881203d58cb2ec0401316e7633738d9110fec570cf76372908be72fb75a50e0754c16fac8a10670d3d7c75c1e210b607020ecdb4b5ce7a7360c1004c79a5d520a54c0592fb8acd2fdca5d735dd53d3ffef837d0f0f2623ca2cebdd587838acf5a0144541336db6ed6b4ab0981b893b91df50a0"),
// 				7426351,
// 				500000000,
// 				100000,
// 				encodingtest.MustDecodeHexZeroX("0xe8468689AB8607fF36663EE6522A7A595Ed8bC0C"),
// 			},
// 			[]byte{0x14, 0xea, 0xc, 0xd0, 0xb4, 0x7, 0xe1, 0x4c, 0xd2, 0x1f, 0x30, 0x7d, 0xf0, 0xed, 0x9f, 0x6d, 0x70, 0xb4, 0xdb, 0x2c, 0xc4, 0x4d, 0xd, 0x1d, 0xec, 0xee, 0xc2, 0xa2, 0xfa, 0x94, 0x7f, 0xbb, 0x20, 0x16, 0x80, 0x48, 0xa1, 0xcd, 0xcc, 0x33, 0x3f, 0x19, 0x96, 0x30, 0x8e, 0xe5, 0xfb, 0xaf, 0x5b, 0xd9, 0xa6, 0x23, 0x60, 0x28, 0xd1, 0x24, 0xbd, 0x15, 0xb5, 0x9e, 0xf0, 0x6e, 0xc, 0x84, 0x1b},
// 			false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := Sender{
// 				key:            tt.fields.key,
// 				networkConfigs: tt.fields.networkConfigs,
// 			}
// 			got, err := e.encode(tt.args.to, tt.args.from, tt.args.data, tt.args.deadline, tt.args.refund, tt.args.gas, tt.args.relayContract)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.encode() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !assert.Equal(t, tt.want, got) {
// 				t.Errorf("Client.encode() = %v, want %v", got, tt.want)
// 			}

// 			hashedData, _ := accounts.TextAndHash(got)

// 			pk := secp256k1test.CharlottePrivateKey
// 			// pk := secp256k1test.SofiaPrivateKey
// 			signed, err := pk.Sign(hashedData)
// 			if err != nil {
// 				panic(err)
// 			}
// 			assert.Equal(t, signed, encodingtest.MustDecodeHexZeroX("0x14ea0cd0b407e14cd21f307df0ed9f6d70b4db2cc44d0d1deceec2a2fa947fbb20168048a1cdcc333f1996308ee5fbaf5bd9a6236028d124bd15b59ef06e0c841b"))

// 		})
// 	}
// }
