package commands

import (
	"context"
	"fmt"
	"time"

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

			addressRPC, _ := cmd.Flags().GetString("rpc-address")
			if addressRPC == "" {
				return errors.New("rpc-address must not be empty")
			}

			rawStorePath, _ := cmd.Flags().GetString("raw-store-path")

			conn, err := newPostgresConnection(cmd)
			if err != nil {
				return err
			}

			defer conn.Close()

			seqProcessor, err := createEthereumProcessor(conn, network, rawStorePath, addressRPC)
			if err != nil {
				return err
			}

			for {
				if err := seqProcessor.NextBlock(context.Background()); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "%+v", err)
				}
				fmt.Println("Infinite Loop 2")
				time.Sleep(time.Second)
			}

			return nil
		},
	}
	cmd.Flags().String("network", ethereum.Mainnet, "Network to run against")
	cmd.Flags().String("rpc-address", "", "Ethereum RPC-JSON address")

	return cmd
}

func createEthereumProcessor(conn *sqlx.DB, network, rawStorePath, addressRPC string) (*processor.Sequential, error) {
	ethClient, err := eth.NewRPC(addressRPC)
	if err != nil {
		return nil, err
	}

	networkID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	syncStore, err := pq.NewSyncStore(conn)
	if err != nil {
		return nil, err
	}

	pubKeyStore, err := pq.NewPublicKeyStore(conn)
	if err != nil {
		return nil, err
	}

	transactionStore, err := pq.NewTransactionStore(conn)
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
		networkID,
	)

	return processor.NewSequential(
		protocols.Ethereum,
		network,
		syncStore,
		eth.NewBlockProcessor(processorTransaction),
		ethClient,
	), nil
}

func sslMode(useSSL bool) string {
	if useSSL {
		return "enable"
	}

	return "disable"
}

// newPostgresConnection returns a connection to a postres database.
// The arguments are parsed from cmd.
func newPostgresConnection(cmd *cobra.Command) (*sqlx.DB, error) {
	host, _ := cmd.Flags().GetString("postgres-host")
	port, _ := cmd.Flags().GetInt("postgres-port")
	useSSL, _ := cmd.Flags().GetBool("postgres-ssl")

	user, err := cmd.Flags().GetString("postgres-user")
	if err != nil {
		return nil, err
	}

	psswd, err := cmd.Flags().GetString("postgres-password")
	if err != nil {
		return nil, err
	}

	dbname, err := cmd.Flags().GetString("postgres-name")
	if err != nil {
		return nil, err
	}

	// use default dbname is not provided
	if dbname == "" {
		dbname = user
	}

	conn, err := pq.NewConnection(user, psswd, dbname, host, sslMode(useSSL), port)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
