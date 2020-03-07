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

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mailchain/mailchain/cmd/internal/http/params"
	"github.com/mailchain/mailchain/crypto"
	ec "github.com/mailchain/mailchain/crypto/cipher/encrypter"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/anysender"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

// SendMessage handler http
func SendMessage(sent stores.Sent, senders map[string]sender.Message, ks keystore.Store, // nolint: funlen
	deriveKeyOptions multi.OptionsBuilders) func(w http.ResponseWriter, r *http.Request) { // nolint: funlen
	// Post swagger:route POST /messages Send SendMessage
	//
	// Send message.
	//
	// Securely send message on the protocol and network specified in the query string to the address.
	// Only the private key holder for the recipient address can decrypt any encrypted contents.
	//
	// - Create mailchain message
	// - Encrypt content with public key
	// - Store message
	// - Encrypt location
	// - Store encrypted location on the blockchain.
	//
	// Responses:
	//   200: StatusOK
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req, err := parsePostRequest(r)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}

		messageSender, ok := senders[fmt.Sprintf("%s/%s", req.Protocol, req.Network)]
		if !ok {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("sender not supported on \"%s/%s\"", req.Protocol, req.Network))
			return
		}

		if messageSender == nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("no sender configured for \"%s/%s\"", req.Protocol, req.Network))
			return
		}

		from, err := address.DecodeByProtocol(req.Body.from.ChainAddress, req.Protocol)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithMessage(err, "failed to decode address"))
			return
		}
		if !ks.HasAddress(from, req.Protocol, req.Network) {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.Errorf("no private key found for `%s` from address", req.Body.Message.Headers.From))
			return
		}

		msg, err := mail.NewMessage(time.Now(), *req.Body.from, *req.Body.to, req.Body.replyTo, req.Body.Message.Subject, []byte(req.Body.Message.Body), req.Body.ContentType)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}

		signer, err := ks.GetSigner(from, req.Protocol, req.Network, signerKind(messageSender), deriveKeyOptions)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(errors.WithMessage(err, "could not get `signer`")))
			return
		}

		encrypter, err := ec.GetEncrypter(req.Body.EncryptionName, req.Body.publicKey)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithMessage(err, "could not get `encrypter`"))
		}

		env, err := envelope.ParseEnvelope(req.Body.Envelope)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithMessage(err, "could not parse `envelope`"))
			return
		}

		if err := mailbox.SendMessage(ctx, req.Protocol, req.Network, msg,
			encrypter, messageSender, sent, signer, env); err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithMessage(err, "could not send message"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func signerKind(s sender.Message) string {
	switch s.(type) {
	case *anysender.Sender:
		return "any.sender"
	default:
		return ""
	}
}

// PostRequest get mailchain inputs
// swagger:parameters SendMessage
type PostRequest struct {
	// Network to use when sending a message.
	//
	// enum: mainnet,goerli,ropsten,rinkeby,local
	// in: query
	// required: true
	// example: goerli
	Network string `json:"network"`

	// Protocol to use when sending a message.
	//
	// enum: ethereum
	// in: query
	// required: true
	// example: ethereum
	Protocol string `json:"protocol"`

	// Message to send
	// in: body
	// required: true
	Body PostRequestBody
}

// parsePostRequest post all the details for the message
func parsePostRequest(r *http.Request) (*PostRequest, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	protocol, err := params.QueryRequireProtocol(r)
	if err != nil {
		return nil, err
	}

	network, err := params.QueryRequireNetwork(r)
	if err != nil {
		return nil, err
	}

	var body PostRequestBody
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.WithMessage(err, "'message' is invalid")
	}

	if err := isValid(&body, protocol, network); err != nil {
		return nil, err
	}

	return &PostRequest{
		Network:  network,
		Protocol: protocol,
		Body:     body,
	}, nil
}

// PostHeaders body
// swagger:model PostMessagesResponseHeaders
type PostHeaders struct {
	// The sender of the message
	// required: true
	// example: Charlotte <5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>
	From string `json:"from"`
	// The recipient of the message
	// required: true
	// To: <4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>
	To string `json:"to"`
	// Reply to if the reply address is different to the from address.
	ReplyTo string `json:"reply-to"`
}

// PostMessage body
// swagger:model PostMessagesResponseMessage
type PostMessage struct {
	// Headers
	// required: true
	// in: body
	Headers *PostHeaders `json:"headers"`
	// Body of the mail message
	// required: true
	// example: Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante,
	// sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu.
	// Integer dictum metus ut sem rutrum aliquet.
	Body string `json:"body"`
	// Subject of the mail message
	// required: true
	// example: Hello world
	Subject string `json:"subject"`
	// Public key of the recipient to encrypt with
	// required: true
	PublicKey string `json:"public-key"`
	// Public key Encoding
	// required: true
	PublicKeyEncoding string `json:"public-key-encoding"`
	// Public key kind
	// require: true
	PublicKeyKind string `json:"public-key-kind"`
}

// PostRequestBody body
// swagger:model SendMessageRequestBody
type PostRequestBody struct {
	// required: true
	Message   PostMessage `json:"message"`
	to        *mail.Address
	from      *mail.Address
	replyTo   *mail.Address
	publicKey crypto.PublicKey
	// Envelope that should be used
	// required: true
	// enum: 0x01, 0x50
	Envelope string `json:"envelope"`
	// Encryption method name
	// required: true
	// enum: aes256cbc, nacl, noop
	EncryptionName string `json:"encryption-method-name"`
	// Message content-type provided by the client
	// required: false (default text/plain; charset=\"UTF-8\")
	// enum: 'text/plain; charset=\"UTF-8\"', 'text/html; charset=\"UTF-8\"'
	ContentType string `json:"content-type"`
}

func checkForEmpties(msg PostMessage) error {
	if msg.Headers == nil {
		return errors.Errorf("headers must not be nil")
	}

	if msg.Body == "" {
		return errors.Errorf("`body` can not be empty")
	}

	if msg.Subject == "" {
		return errors.Errorf("`subject` can not be empty")
	}

	if msg.PublicKey == "" {
		return errors.Errorf("`public-key` can not be empty")
	}

	if msg.PublicKeyEncoding == "" {
		return errors.Errorf("`public-key-encoding` can not be empty")
	}

	return nil
}

func isValid(p *PostRequestBody, protocol, network string) error {
	if p == nil {
		return errors.New("PostRequestBody must not be nil")
	}

	if err := checkForEmpties(p.Message); err != nil {
		return err
	}

	var err error

	p.to, err = mail.ParseAddress(p.Message.Headers.To, protocol, network)
	if err != nil {
		return errors.WithMessage(err, "`to` is invalid")
	}

	p.from, err = mail.ParseAddress(p.Message.Headers.From, protocol, network)
	if err != nil {
		return errors.WithMessage(err, "`from` is invalid")
	}

	if p.Message.Headers.ReplyTo != "" {
		p.replyTo, err = mail.ParseAddress(p.Message.Headers.ReplyTo, protocol, network)
		if err != nil {
			return errors.WithMessage(err, "`reply-to` is invalid")
		}
	}

	// Validate Public-key-encoding
	if _, validEncod := encoding.PublicKeyEncoding()[p.Message.PublicKeyEncoding]; !validEncod {
		return errors.Errorf("invalid `public-key-encoding` ")
	}

	// Validate Public-key Kind
	if _, validkindType := crypto.KeyTypes()[p.Message.PublicKeyKind]; !validkindType {
		return errors.Errorf("invalid `public-key-kind` ")
	}

	keyBytes, err := encoding.DecodeHexZeroX(p.Message.PublicKey)
	if err != nil {
		return errors.WithMessage(err, "invalid `public-key-bytes`")
	}

	p.publicKey, err = multikey.PublicKeyFromBytes(p.Message.PublicKeyKind, keyBytes)
	if err != nil {
		return errors.WithMessage(err, "invalid `public-key`")
	}

	if p.Envelope == "" {
		return errors.Errorf("`envelope` can not be empty")
	}

	if p.EncryptionName == "" {
		return errors.Errorf("`encryption-method-name` can not be empty")
	}

	return nil
}
