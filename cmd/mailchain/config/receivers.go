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

// nolint: dupl
package config

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

func SetReceiver(vpr *viper.Viper, chain, network, receiver string) error {
	viper.Set(fmt.Sprintf("chains.%s.networks.%s.receiver", chain, network), receiver)
	if err := setClient(vpr, receiver, network); err != nil {
		return err
	}
	fmt.Printf("%s used for receiving messages\n", receiver)
	return nil
}

// GetReceivers in configured state
func GetReceivers(vpr *viper.Viper) (map[string]mailbox.Receiver, error) {
	receivers := make(map[string]mailbox.Receiver)
	for chain := range viper.GetStringMap("chains") {
		chRcvrs, err := getChainReceivers(vpr, chain)
		if err != nil {
			return nil, err
		}
		if err := mergo.Merge(&receivers, chRcvrs); err != nil {
			return nil, err
		}
	}
	return receivers, nil
}

func getChainReceivers(vpr *viper.Viper, chain string) (map[string]mailbox.Receiver, error) {
	receivers := make(map[string]mailbox.Receiver)
	for network := range viper.GetStringMap(fmt.Sprintf("chains.%s.networks", chain)) {
		receiver, err := getReceiver(vpr, chain, network)
		if err != nil {
			return nil, err
		}
		receivers[fmt.Sprintf("%s.%s", chain, network)] = receiver
	}

	return receivers, nil
}

func getReceiver(vpr *viper.Viper, chain, network string) (mailbox.Receiver, error) {
	switch viper.GetString(fmt.Sprintf("chains.%s.networks.%s.receiver", chain, network)) {
	case names.Etherscan:
		return getEtherscanClient(vpr)
	case names.EtherscanNoAuth:
		return getEtherscanNoAuthClient()
	default:
		return nil, errors.Errorf("unsupported receiver")
	}
}
