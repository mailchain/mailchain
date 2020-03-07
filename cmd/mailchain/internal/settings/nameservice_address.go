package settings //nolint: dupl

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/pkg/errors"
)

func addressNameServices(s values.Store) *AddressNameServices {
	return &AddressNameServices{
		clients: map[string]NameServiceAddressClient{
			defaults.Mailchain: mailchainAddressNameServices(s),
		},
	}
}

// AddressNameServices configuration element.
type AddressNameServices struct {
	clients map[string]NameServiceAddressClient
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s AddressNameServices) Output() output.Element {
	elements := []output.Element{}
	for _, c := range s.clients {
		elements = append(elements, c.Output())
	}
	return output.Element{
		FullName: "nameservice-address",
		Elements: elements,
	}
}

// Produce `nameservice.ReverseLookup` based on configuration settings.
func (s AddressNameServices) Produce(client string) (nameservice.ReverseLookup, error) {
	if client == "" {
		return nil, nil
	}
	m, ok := s.clients[client]
	if !ok {
		return nil, errors.Errorf("%q not a supported address name service", client)
	}
	return m.Produce()
}

func mailchainAddressNameServices(s values.Store) *MailchainAddressNameServices {
	enabledNetworks := []string{}
	for _, n := range ethereum.Networks() {
		enabledNetworks = append(enabledNetworks, "ethereum/"+n)
	}
	return &MailchainAddressNameServices{
		BaseURL: values.NewDefaultString("https://ns.mailchain.xyz/", s, "nameservice-address.mailchain.base-url"),
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"nameservice-address.mailchain.enabled-networks",
		),
	}
}

// MailchainAddressNameServices configuration element.
type MailchainAddressNameServices struct {
	BaseURL                 values.String
	EnabledProtocolNetworks values.StringSlice
}

// Produce `nameservice.ReverseLookup` based on configuration settings.
func (s MailchainAddressNameServices) Produce() (nameservice.ReverseLookup, error) {
	return nameservice.NewLookupService(s.BaseURL.Get()), nil
}

// Supports a map of what protocol and network combinations are supported.
func (s MailchainAddressNameServices) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range s.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s MailchainAddressNameServices) Output() output.Element {
	return output.Element{
		FullName: "nameservice-address.mailchain",
		Attributes: []output.Attribute{
			s.BaseURL.Attribute(),
			s.EnabledProtocolNetworks.Attribute(),
		},
	}
}
