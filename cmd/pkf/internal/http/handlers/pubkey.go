package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/cmd/internal/http/params"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/pubkey"
	"github.com/pkg/errors"
)

// GetPublicKey returns a handler get spec
func GetPublicKey(store datastore.PublicKeyStore) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /public-key PublicKey GetPublicKey
	//
	// Public key from address.
	//
	// This method will get the public key to use when encrypting messages and envelopes.
	// Protocols and networks have different methods for retrieving or calculating a public key from an address.
	//
	// Responses:
	//   200: GetPublicKeyResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := parseGetPublicKey(r)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}

		publicKey, err := store.GetPublicKey(ctx, req.Protocol, req.Network, req.addressBytes)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		encodedKey, encodingType, err := pubkey.EncodeByProtocol(publicKey.PublicKey.Bytes(), req.Protocol)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		encryptionTypes, err := pubkey.EncryptionMethods(publicKey.PublicKey.Kind())
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetPublicKeyResponseBody{
			PublicKey:                encodedKey,
			PublicKeyEncoding:        encodingType,
			SupportedEncryptionTypes: encryptionTypes,
			PublicKeyKind:            publicKey.PublicKey.Kind(),
			BlockHash:                encoding.EncodeHexZeroX(publicKey.BlockHash),
			TxHash:                   encoding.EncodeHexZeroX(publicKey.TxHash),
		})
	}
}

// parseGetPublicKey get all the details for the get request
func parseGetPublicKey(r *http.Request) (*GetPublicKeyRequest, error) {
	protocol, err := params.QueryRequireProtocol(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	network, err := params.QueryRequireNetwork(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	addr, err := params.QueryRequireAddress(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	addressBytes, err := address.DecodeByProtocol(addr, protocol)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to decode public key")
	}

	return &GetPublicKeyRequest{
		Address:      addr,
		addressBytes: addressBytes,
		Network:      network,
		Protocol:     protocol,
	}, nil
}

// GetPublicKeyRequest pubic key from address request
// swagger:parameters GetPublicKey
type GetPublicKeyRequest struct {
	// Address to to use when performing public key lookup.
	//
	// in: query
	// required: true
	// example: 0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae
	// pattern: 0x[a-fA-F0-9]{40}
	Address      string `json:"address"`
	addressBytes []byte

	// Network to use when performing public key lookup.
	//
	// enum: mainnet,goerli,ropsten,rinkeby,local
	// in: query
	// required: true
	// example: goerli
	Network string `json:"network"`

	// Protocol to use when performing public key lookup.
	//
	// enum: ethereum
	// in: query
	// required: true
	// example: ethereum
	Protocol string `json:"protocol"`
}

// GetPublicKeyResponse public key from address response
//
// swagger:response GetPublicKeyResponse
type GetPublicKeyResponse struct {
	// in: body
	Body GetPublicKeyResponseBody
}

// GetPublicKeyResponseBody body response
//
// swagger:model GetPublicKeyResponseBody
type GetPublicKeyResponseBody struct {
	// The public key encoded as per `public_key_encoding`
	//
	// Required: true
	// example: 0x79964e63752465973b6b3c610d8ac773fc7ce04f5d1ba599ba8768fb44cef525176f81d3c7603d5a2e466bc96da7b2443bef01b78059a98f45d5c440ca379463
	PublicKey string `json:"public_key"`

	// Encoding method used for encoding the `public_key`
	//
	// Required: true
	// example: hex/0x-prefix
	PublicKeyEncoding string `json:"public_key_encoding"`

	// Encoding method used for encoding the `public_key`
	//
	// Required: true
	// example: ["secp256k1", "sr25519", "ed25519"]
	PublicKeyKind string `json:"public_key_kind"`

	// Supported encryption methods for public keys.
	//
	// Required: true
	// example: ["aes256cbc", "nacl-ecdh", "noop"]
	SupportedEncryptionTypes []string `json:"supported_encryption_types"`

	// Block hash where we found the PublicKey
	//
	// Required: true
	// example: 0x373d339e45a701447367d7b9c7cef84aab79c2b2714271b908cda0ab3ad0849b
	BlockHash string `json:"block_hash"`

	// Transaction hash where we found the PublicKey
	//
	// Required: true
	// example: 0x98beb27135aa0a25650557005ad962919d6a278c4b3dde7f4f6a3a1e65aa746c
	TxHash string `json:"tx_hash"`
}