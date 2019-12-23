package encrypter

import (
	keys "github.com/mailchain/mailchain/crypto"
	crypto "github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/cipher/nacl"
	"github.com/pkg/errors"
)

// Cipher Name lookup
const (
	// NoOperation encryption type name.
	NoOperation string = "noop"
	// NACL encryption type name.
	NACL string = "nacl"
	// AES256CBC encryption type name.
	AES256CBC string = "aes256cbc"
)

// GetEncrypter is an `Encrypter` factory that returns an encrypter
func GetEncrypter(encryption string, pubKey keys.PublicKey) (crypto.Encrypter, error) {
	switch encryption {
	case AES256CBC:
		return aes256cbc.NewEncrypter(pubKey)
	case NACL:
		return nacl.NewEncrypter(pubKey)
	case "":
		return nil, errors.Errorf("`encryption` provided is set to empty")
	default:
		return nil, errors.Errorf("`encryption` provided is invalid")
	}
}
