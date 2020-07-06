// Copyright 2020 Finobo
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

package blockscout

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// NewAPIClient create new API client
func NewAPIClient() (*APIClient, error) {
	return &APIClient{
		networkConfigs: map[string]networkConfig{
			ethereum.Mainnet: {
				url:    "https://blockscout.com/eth/mainnet/api",
				rpcURL: "https://relay.mailchain.xyz/json-rpc/ethereum/mainnet",
			},
		},
	}, nil
}

// APIClient for talking to etherscan
type APIClient struct {
	networkConfigs map[string]networkConfig
}

type networkConfig struct {
	rpcURL string
	url    string
}

// IsNetworkSupported checks if the network is supported by etherscan API
func (c APIClient) isNetworkSupported(network string) bool {
	_, ok := c.networkConfigs[network]
	return ok
}

// GetTransactionsByAddress get transactions from address via etherscan
func (c APIClient) getTransactionsByAddress(network string, address []byte) (*txList, error) {
	config, ok := c.networkConfigs[network]
	if !ok {
		return nil, errors.Errorf("network not supported")
	}

	txListResponse, err := resty.R().
		SetQueryParams(map[string]string{
			"module":     "account",
			"action":     "txlist",
			"startblock": "0",
			"endblock":   "99999999",
			"sort":       "desc",
			"address":    common.BytesToAddress(address).Hex(),
		}).Get(config.url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if txListResponse.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("status: %d body: %s", txListResponse.StatusCode(), string(txListResponse.Body()))
	}

	txResult := &txList{}
	if err := json.Unmarshal(txListResponse.Body(), txResult); err != nil {
		return nil, errors.WithMessage(err, string(txListResponse.Body()))
	}

	return txResult, nil
}

// GetTransactionByHash get transaction details from transaction hash via etherscan
func (c APIClient) getTransactionByHash(network string, hash common.Hash) (*types.Transaction, error) {
	config, ok := c.networkConfigs[network]
	if !ok {
		return nil, errors.Errorf("network not supported")
	}

	client, err := ethclient.Dial(config.rpcURL)
	if err != nil {
		return nil, err
	}

	tx, _, err := client.TransactionByHash(context.Background(), hash)

	return tx, err
}
