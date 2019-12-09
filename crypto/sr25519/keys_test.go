package sr25519

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/internal/encoding/encodingtest"
)

var ( //nolint
	sofiaSeed       = encodingtest.MustDecodeHex("5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a") //nolint: lll
	sofiaPrivateKey = PrivateKey{key: func() *schnorrkel.SecretKey {
		priv, err := keyFromSeed([]byte{0x5c, 0x6d, 0x7a, 0xdf, 0x75, 0xbd, 0xa1, 0x18, 0xc, 0x22, 0x5d, 0x25, 0xf3, 0xaa, 0x8d, 0xc1, 0x74, 0xbb, 0xfb, 0x3c, 0xdd, 0xee, 0x11, 0xae, 0x9a, 0x85, 0x98, 0x2f, 0x6f, 0xaf, 0x79, 0x1a})
		if err != nil {
			panic(err)
		}
		return priv
	}(),
	}

	sofiaPrivateKeyBytes = []byte{0x5c, 0x6d, 0x7a, 0xdf, 0x75, 0xbd, 0xa1, 0x18, 0xc, 0x22, 0x5d, 0x25, 0xf3, 0xaa, 0x8d, 0xc1, 0x74, 0xbb, 0xfb, 0x3c, 0xdd, 0xee, 0x11, 0xae, 0x9a, 0x85, 0x98, 0x2f, 0x6f, 0xaf, 0x79, 0x1a}
	sofiaPublicKey       = newPublicKey([]byte{0x16, 0x9a, 0x11, 0x72, 0x18, 0x51, 0xf5, 0xdf, 0xf3, 0x54, 0x1d, 0xd5, 0xc4, 0xb0, 0xb4, 0x78, 0xac, 0x1c, 0xd0, 0x92, 0xc9, 0xd5, 0x97, 0x6e, 0x83, 0xda, 0xa0, 0xd0, 0x3f, 0x26, 0x62, 0xc}) //nolint: lll
	sofiaPublicKeyBytes  = []byte{0x16, 0x9a, 0x11, 0x72, 0x18, 0x51, 0xf5, 0xdf, 0xf3, 0x54, 0x1d, 0xd5, 0xc4, 0xb0, 0xb4, 0x78, 0xac, 0x1c, 0xd0, 0x92, 0xc9, 0xd5, 0x97, 0x6e, 0x83, 0xda, 0xa0, 0xd0, 0x3f, 0x26, 0x62, 0xc}               //nolint: lll
	charlotteSeed        = encodingtest.MustDecodeHex("23b063a581fd8e5e847c4e2b9c494247298791530f5293be369e8bf23a45d2bd")                                                                                                                      //nolint: lll
	charlottePrivateKey  = PrivateKey{key: func() *schnorrkel.SecretKey {
		priv, err := keyFromSeed([]byte{0x23, 0xb0, 0x63, 0xa5, 0x81, 0xfd, 0x8e, 0x5e, 0x84, 0x7c, 0x4e, 0x2b, 0x9c, 0x49, 0x42, 0x47, 0x29, 0x87, 0x91, 0x53, 0xf, 0x52, 0x93, 0xbe, 0x36, 0x9e, 0x8b, 0xf2, 0x3a, 0x45, 0xd2, 0xbd})
		if err != nil {
			panic(err)
		}

		return priv
	}()}
	charlottePrivateKeyBytes = []byte{0x23, 0xb0, 0x63, 0xa5, 0x81, 0xfd, 0x8e, 0x5e, 0x84, 0x7c, 0x4e, 0x2b, 0x9c, 0x49, 0x42, 0x47, 0x29, 0x87, 0x91, 0x53, 0xf, 0x52, 0x93, 0xbe, 0x36, 0x9e, 0x8b, 0xf2, 0x3a, 0x45, 0xd2, 0xbd}
	charlottePublicKey       = newPublicKey([]byte{0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30}) //nolint: lll
	charlottePublicKeyBytes  = []byte{0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30}               //nolint: lll
) //nolint: lll
