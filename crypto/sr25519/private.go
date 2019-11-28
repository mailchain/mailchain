package sr25519

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/crypto"

	"github.com/pkg/errors"
)

const (
	keyPairSize    = 96
	privateKeySize = 64
	seedSize       = 32
)

var SigningContext = []byte("substrate") //nolint gochecknoglobals

// Private Key sr25519
type PrivateKey struct {
	key *schnorrkel.SecretKey
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	b := pk.key.Encode()
	kb := make([]byte, len(b))
	copy(kb, b[:])

	return kb
}

// Kind is the type of private key.
func (pk PrivateKey) Kind() string {
	return crypto.SR25519
}

// PublicKey return the public key that is derived from the private key
func (pk PrivateKey) PublicKey() crypto.PublicKey {
	pub, err := pk.key.Public()
	if err != nil {
		return nil
	}

	return &PublicKey{key: pub}
}

// Sign uses the private key to sign the message using the sr25519 signature algorithm
func (pk PrivateKey) Sign(message []byte) (signature []byte, err error) {
	if pk.key == nil {
		return nil, errors.New("key is nil")
	}

	t := schnorrkel.NewSigningContext(SigningContext, message)

	sig, err := pk.key.Sign(t)
	if err != nil {
		return nil, err
	}

	enc := sig.Encode()

	return enc[:], nil
}

// Encode returns the 32-byte encoding of the private key
func (pk *PrivateKey) Encode() []byte {
	if pk.key == nil {
		return nil
	}

	enc := pk.key.Encode()

	return enc[:]
}

// Decode decodes the input bytes into a private key and sets the receiver the decoded key
// Input must be 32 bytes, or else this function will error
func (pk *PrivateKey) Decode(in []byte) error {
	if len(in) != privateKeySize {
		return errors.New("input to sr25519 private key decode is not 32 bytes")
	}

	b := [32]byte{}
	copy(b[:], in)

	pk.key = &schnorrkel.SecretKey{}

	return pk.key.Decode(b)
}

func keyFromSeed(b []byte) (*schnorrkel.SecretKey, error) {
	kb := [32]byte{}
	copy(b, kb[:])

	priv, err := schnorrkel.NewMiniSecretKeyFromRaw(kb)
	if err != nil {
		return nil, err
	}

	return priv.ExpandUniform(), nil
}

// PrivateKeyFromBytes get a private key from seed []byte
func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	switch len(privKey) {
	case privateKeySize:
		privKey, err := keyFromSeed(privKey)
		if err != nil {
			return nil, err
		}

		return &PrivateKey{key: privKey}, nil
	case seedSize:
		privKey, err := keyFromSeed(privKey)
		if err != nil {
			return nil, err
		}

		return &PrivateKey{key: privKey}, nil
	case keyPairSize:
		privKey, err := keyFromSeed(privKey)
		if err != nil {
			return nil, err
		}

		pk, err := NewKeypair(privKey)
		if err != nil {
			return nil, err
		}

		return pk.private, nil
	default:
		return nil, errors.Errorf("bad key length")
	}
}