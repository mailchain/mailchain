package anysender

import (
	"context"
	"fmt"
	"math/big"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/sender"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
	"gopkg.in/resty.v1"
)

// NewSender create new API client
func NewSender(address string, relayContractAddress, rpcURL string) (*Sender, error) {
	return &Sender{
		address:              address,
		relayContractAddress: relayContractAddress,
		rpcURL:               rpcURL,
	}, nil
}

// Sender for talking to any.sender
type Sender struct {
	address              string
	relayContractAddress string
	rpcURL               string
}

func (e Sender) Send(ctx context.Context, network string, to, from, data []byte, txSigner signer.Signer, opts sender.SendOpts) error {
	ethClient, err := ethclient.Dial(e.rpcURL)
	if err != nil {
		return err
	}
	currentBlock, err := ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return err
	}
	toAddress := common.BytesToAddress(to)

	gasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return errors.WithMessage(err, "could not determine gas price")
	}

	gas, err := ethClient.EstimateGas(ctx, geth.CallMsg{
		Data:     data,
		From:     common.BytesToAddress(from),
		GasPrice: gasPrice,
		To:       &toAddress,
		Value:    big.NewInt(0),
	})
	if err != nil {
		return err
	}

	relayContractAddress, err := encoding.DecodeHexZeroX(e.relayContractAddress)
	if err != nil {
		return err
	}
	deadline := currentBlock.Number().Int64() + 500 // this is an arbitrary number
	refund := int64(500000000)

	finalGas := gas + 60000

	signedTransaction, err := txSigner.Sign(SignerOptions{
		from:          from,
		to:            to,
		data:          data,
		deadline:      deadline,
		refund:        refund,
		gas:           int64(finalGas),
		relayContract: relayContractAddress,
	})
	if err != nil {
		return err
	}

	signature, ok := signedTransaction.([]byte)
	if !ok {
		return errors.Errorf("unknown signature type")
	}

	req := resty.R().
		SetBody(map[string]interface{}{
			"from":                 encoding.EncodeHexZeroX(from),
			"to":                   encoding.EncodeHexZeroX(to),
			"data":                 encoding.EncodeHexZeroX(data),
			"deadlineBlockNumber":  deadline,
			"gas":                  finalGas,
			"refund":               fmt.Sprintf("%d", refund),
			"relayContractAddress": encoding.EncodeHexZeroX(relayContractAddress),
			"signature":            encoding.EncodeHexZeroX(signature),
		}).
		// SetBody(sendReq{
		// 	From:                 encoding.EncodeHexZeroX(from),
		// 	To:                   encoding.EncodeHexZeroX(to),
		// 	Data:                 encoding.EncodeHexZeroX(data),
		// 	DeadlineBlockNumber:  deadline,
		// 	Gas:                  gasPrice.Int64(),
		// 	Refund:               string(refund),
		// 	RelayContractAddress: encoding.EncodeHexZeroX(relayContractAddress),
		// 	Signature:            encoding.EncodeHexZeroX(signature),
		// })
		SetHeader("Content-Type", "application/json")

	response, err := req.Post(e.address + "/relay")
	if err != nil {
		println(fmt.Sprintf("%+v", req.Body))
		return errors.WithStack(err)
	}
	if response.StatusCode() != 200 {
		return errors.Errorf("request: %+v\nresponse: %+v\n", req.Body, response.String())
		// return errors.New(response.String())
	}

	return nil
}

func (e Sender) encode(to []byte, from []byte, data []byte, deadline int64, refund int64, gas int64, relayContract []byte) ([]byte, error) {
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

	packed, err := abi.Arguments{
		{Type: typeAddress}, // to
		{Type: typeAddress}, // from
		{Type: typeBytes},   // data
		{Type: typeUnit256}, // deadline
		{Type: typeUnit256}, // refund
		{Type: typeUnit256}, // gas
		{Type: typeAddress}, // relayContract
	}.Pack(
		common.BytesToAddress(to),
		common.BytesToAddress(from),
		data,
		big.NewInt(deadline),
		big.NewInt(refund),
		big.NewInt(gas),
		common.BytesToAddress(relayContract),
	)
	if err != nil {
		return nil, err
	}

	var buf []byte
	hash := sha3.NewLegacyKeccak256()
	_, _ = hash.Write(packed)
	buf = hash.Sum(buf)

	return buf, nil
}
