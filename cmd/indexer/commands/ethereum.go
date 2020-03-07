package commands

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/params"
	"github.com/jmoiron/sqlx"
	eth "github.com/mailchain/mailchain/cmd/indexer/internal/ethereum"
	"github.com/mailchain/mailchain/cmd/indexer/internal/processor"
	"github.com/mailchain/mailchain/cmd/internal/datastore/os"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ethereumCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "ethereum",
		Short:            "run ethereum sequential processor",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			network, _ := cmd.Flags().GetString("network")
			protocol, _ := cmd.Flags().GetString("protocol")
			blockNumber, _ := cmd.Flags().GetUint64("start-block")

			addressRPC, _ := cmd.Flags().GetString("rpc-address")
			if addressRPC == "" {
				return errors.New("rpc-address must not be empty")
			}

			rawStorePath, _ := cmd.Flags().GetString("raw-store-path")

			connIndexer, err := newPostgresConnection(cmd, "indexer")
			if err != nil {
				return err
			}

			connPublicKey, err := newPostgresConnection(cmd, "pubkey")
			if err != nil {
				return err
			}

			connEnvelope, err := newPostgresConnection(cmd, "envelope")
			if err != nil {
				return err
			}

			defer connIndexer.Close()
			defer connPublicKey.Close()
			defer connEnvelope.Close()

			seqProcessor, err := createEthereumProcessor(connIndexer, connPublicKey, connEnvelope, blockNumber, protocol, network, rawStorePath, addressRPC)
			if err != nil {
				return err
			}

			for {
				err := seqProcessor.NextBlock(context.Background())

				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "%+v", err)
				}
			}
		},
	}

	cmd.Flags().Uint64("start-block", 0, "Block number from which the indexer will start")
	cmd.Flags().String("protocol", protocols.Ethereum, "Protocol to run against")
	cmd.Flags().String("network", ethereum.Mainnet, "Network to run against")
	cmd.Flags().String("rpc-address", "", "Ethereum RPC-JSON address")

	return cmd
}

func createEthereumProcessor(connIndexer, connPublicKey, connEnvelope *sqlx.DB, blockNumber uint64, protocol, network, rawStorePath, addressRPC string) (*processor.Sequential, error) {
	ctx := context.Background()

	ethClient, err := eth.NewRPC(addressRPC)
	if err != nil {
		return nil, err
	}

	networkID, err := ethClient.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	chainCfg, err := chainConfig(network)
	if err != nil {
		return nil, err
	}

	if chainCfg.ChainID.Cmp(networkID) != 0 {
		return nil, errors.Errorf("networkID from RPC does not match chain config network ID")
	}

	syncStore, err := pq.NewSyncStore(connIndexer)
	if err != nil {
		return nil, err
	}

	pubKeyStore, err := pq.NewPublicKeyStore(connPublicKey)
	if err != nil {
		return nil, err
	}

	transactionStore, err := pq.NewTransactionStore(connEnvelope)
	if err != nil {
		return nil, err
	}

	rawStore, err := os.NewRawTransactionStore(rawStorePath)
	if err != nil {
		return nil, err
	}

	processorTransaction := eth.NewTransactionProcessor(
		transactionStore,
		rawStore,
		pubKeyStore,
		chainCfg,
	)

	if err := syncStore.PutBlockNumber(ctx, protocol, network, blockNumber); err != nil {
		return nil, err
	}

	return processor.NewSequential(
		protocols.Ethereum,
		network,
		syncStore,
		eth.NewBlockProcessor(processorTransaction),
		ethClient,
	), nil
}

func chainConfig(network string) (*params.ChainConfig, error) {
	switch network {
	case ethereum.Goerli:
		return params.GoerliChainConfig, nil
	case ethereum.Mainnet:
		return params.MainnetChainConfig, nil
	case ethereum.Rinkeby:
		return params.RinkebyChainConfig, nil
	case ethereum.Ropsten:
		return params.TestnetChainConfig, nil
	default:
		return nil, errors.Errorf("can not determine chain config from network: %s", network)
	}
}
