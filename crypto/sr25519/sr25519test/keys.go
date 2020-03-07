package sr25519test

import (
	"log"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/mailchain/mailchain/encoding/encodingtest"
)

// SofiaPrivateKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var SofiaPrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// SofiaPublicKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var SofiaPublicKey crypto.PublicKey //nolint: gochecknoglobals
// CharlottePrivateKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var CharlottePrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// CharlottePublicKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var CharlottePublicKey crypto.PublicKey //nolint: gochecknoglobals
// EvePrivateKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var EvePrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// EvePublicKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var EvePublicKey crypto.PublicKey //nolint: gochecknoglobals

//nolint: gochecknoinits
func init() {
	var err error
	SofiaPrivateKey, err = sr25519.PrivateKeyFromBytes(encodingtest.MustDecodeHex("5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a")) //nolint: lll
	if err != nil {
		log.Fatal(err)
	}

	SofiaPublicKey = SofiaPrivateKey.PublicKey()

	CharlottePrivateKey, err = sr25519.PrivateKeyFromBytes(encodingtest.MustDecodeHex("23b063a581fd8e5e847c4e2b9c494247298791530f5293be369e8bf23a45d2bd")) //nolint: lll
	if err != nil {
		log.Fatal(err)
	}

	CharlottePublicKey = CharlottePrivateKey.PublicKey()

	EvePrivateKey, err = sr25519.PrivateKeyFromBytes(encodingtest.MustDecodeHex("000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f")) //nolint: lll
	if err != nil {
		log.Fatal(err)
	}

	EvePublicKey = EvePrivateKey.PublicKey()
}
