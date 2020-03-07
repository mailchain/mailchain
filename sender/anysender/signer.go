package anysender

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

// SignerOptions options that can be set when signing an ethereum transaction.
type SignerOptions struct {
	to            []byte
	from          []byte
	data          []byte
	deadline      int64
	refund        int64
	gas           int64
	relayContract []byte
	// Tx      *types.Transaction
	// ChainID *big.Int
}

// NewSigner returns a new ethereum signer that can be used to sign transactions.
func NewSigner(privateKey crypto.PrivateKey) (*Signer, error) {
	if err := validatePrivateKeyType(privateKey); err != nil {
		return nil, errors.WithStack(err)
	}

	return &Signer{privateKey: privateKey}, nil
}

// Sign an ethereum transaction with the private key.
func (e Signer) Sign(opts signer.Options) (signedTransaction interface{}, err error) {
	if opts == nil {
		return nil, errors.New("opts must not be nil")
	}

	if err := validatePrivateKeyType(e.privateKey); err != nil {
		return nil, err
	}

	switch opts := opts.(type) {
	case SignerOptions:
		encoded, err := e.encodeABI(opts.to, opts.from, opts.data, opts.deadline, opts.refund, opts.gas, opts.relayContract)
		if err != nil {
			return nil, err
		}

		var hashed []byte
		hash := sha3.NewLegacyKeccak256()
		_, _ = hash.Write(encoded)

		hashedWithText, _ := accounts.TextAndHash(hash.Sum(hashed))

		signature, err := e.privateKey.Sign(hashedWithText)
		if err != nil {
			return nil, err
		}

		v := signature[64]
		if v != 27 && v != 28 {
			v = 27 + (v % 2)
		}

		if v != 0x1c {
			v = 0x1b
		}

		signature[64] = v

		return signature, nil
	default:
		return nil, errors.New("invalid options for any.sender signing")
	}
}

// Signer for ethereum.
type Signer struct {
	privateKey crypto.PrivateKey
}

func validatePrivateKeyType(pk crypto.PrivateKey) error {
	switch pk.(type) {
	case secp256k1.PrivateKey:
		return nil
	case *secp256k1.PrivateKey:
		return nil
	default:
		return errors.New("invalid key type")
	}
}

func (e Signer) encodeABI(to []byte, from []byte, data []byte, deadline int64, refund int64, gas int64, relayContract []byte) ([]byte, error) {
	typeAddress, err := abi.NewType("address", "address", nil)
	if err != nil {
		return nil, err
	}

	typeBytes, err := abi.NewType("bytes", "bytes", nil)
	if err != nil {
		return nil, err
	}

	typeUnit256, err := abi.NewType("uint256", "uint256", nil)
	if err != nil {
		return nil, err
	}

	args := abi.Arguments{
		{Type: typeAddress}, // to
		{Type: typeAddress}, // from
		{Type: typeBytes},   // data
		{Type: typeUnit256}, // deadline
		{Type: typeUnit256}, // refund
		{Type: typeUnit256}, // gas
		{Type: typeAddress}, // relayContract
	}

	return args.Pack(
		common.BytesToAddress(to),
		common.BytesToAddress(from),
		data,
		big.NewInt(deadline),
		big.NewInt(refund),
		big.NewInt(gas),
		common.BytesToAddress(relayContract),
	)
}
