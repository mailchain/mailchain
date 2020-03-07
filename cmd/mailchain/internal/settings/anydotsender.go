package settings

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/anysender"
)

func anyDotSenderSender(s values.Store, relayContractAddress string, network string) *AnyDotSender {
	return &AnyDotSender{network: network,
		Address: values.NewDefaultString(fmt.Sprintf("https://api.pisa.watch/any.sender.%s", network), s,
			fmt.Sprintf("senders.ethereum-anydotsender-%s.address", network)),
		RPCURL: values.NewDefaultString(fmt.Sprintf("https://relay.mailchain.xyz/json-rpc/ethereum/%s", network), s,
			fmt.Sprintf("senders.ethereum-anydotsender-%s.rpc-url", network)),
		RelayContractAddress: values.NewDefaultString(relayContractAddress, s,
			fmt.Sprintf("senders.ethereum-anydotsender-%s.relay-contract-address", relayContractAddress)),
	}
}

// AnyDotSender configuration element.
type AnyDotSender struct {
	Address              values.String
	RelayContractAddress values.String
	RPCURL               values.String
	network              string
}

// Produce `sender.Message` based on configuration settings.
func (e AnyDotSender) Produce() (sender.Message, error) {
	return anysender.NewSender(e.Address.Get(), e.RelayContractAddress.Get(), e.RPCURL.Get())
}

// Supports a map of what protocol and network combinations are supported.
func (e AnyDotSender) Supports() map[string]bool {
	return map[string]bool{
		"ethereum/" + e.network: true,
	}
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (e AnyDotSender) Output() output.Element {
	return output.Element{
		FullName: "senders.ethereum-anydotsender-" + e.network,
		Attributes: []output.Attribute{
			e.Address.Attribute(),
			e.RPCURL.Attribute(),
			e.RelayContractAddress.Attribute(),
		},
	}
}
