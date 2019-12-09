package sr25519

import (
	"errors"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/crypto"
)

const (
	publicKeySize   = 32
	signatureLength = 64
	signatureSize   = 64
)

// PublicKey is a interface
type PublicKey struct {
	key *schnorrkel.PublicKey
}

// Bytes return Publickey Bytes
func (pk PublicKey) Bytes() []byte {
	b := pk.key.Encode()

	return b[:]
}

// Kind returns the key type
func (pk PublicKey) Kind() string {
	return crypto.SR25519
}

func newPublicKey(b []byte) PublicKey { //nolint
	enc := [publicKeySize]byte{}
	copy(enc[:], b)

	public := schnorrkel.NewPublicKey(enc)

	return PublicKey{key: public}
}

// Verify uses the sr25519 signature algorithm to verify that the message was signed by
// this public key; it returns true if this key created the signature for the message,
// false otherwise
func (pk PublicKey) Verify(message, sig []byte) bool {
	b := [signatureLength]byte{}
	copy(b[:], sig)

	s := &schnorrkel.Signature{}

	err := s.Decode(b)
	if err != nil {
		return false
	}

	t := schnorrkel.NewSigningContext(SigningContext, message)

	return pk.key.Verify(s, t)
}

// Encode returns the 32-byte encoding of the public key
func (pk PublicKey) Encode() []byte {
	enc := pk.key.Encode()

	return enc[:]
}

// Decode decodes the input bytes into a public key and sets the receiver the decoded key
// Input must be 32 bytes, or else this function will error
func (pk PublicKey) Decode(in []byte) error {
	b := [publicKeySize]byte{}
	copy(b[:], in)

	return pk.key.Decode(b)
}

func schnorrkelPublicKeyFromBytes(in []byte) (*schnorrkel.PublicKey, error) {
	if len(in) != publicKeySize {
		return nil, errors.New("input to sr25519 public key decode is not 32 bytes")
	}

	b := [32]byte{}
	copy(b[:], in)

	key := &schnorrkel.PublicKey{}
	err := key.Decode(b)

	return key, err
}

// PublicKeyFromBytes - Convert byte array to PublicKey
func PublicKeyFromBytes(keyBytes []byte) (*PublicKey, error) {
	switch len(keyBytes) {
	case publicKeySize:
		pubKey, _ := schnorrkelPublicKeyFromBytes(keyBytes)

		return &PublicKey{pubKey}, nil
	default:
		return nil, errors.New("public key must be 32 bytes")
	}
}