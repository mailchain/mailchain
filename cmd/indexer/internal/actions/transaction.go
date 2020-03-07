package actions

import (
	"bytes"
	"context"

	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/encoding"
)

//go:generate mockgen -source=transaction.go -package=actionstest -destination=./actionstest/transaction_mock.go

// Transaction processes transactions
type Transaction interface {
	// Run the transaction processor
	Run(ctx context.Context, protocol, network string, tx interface{}, txOpts TransactionOptions) error
}

// TransactionOptions related to different transactions
type TransactionOptions interface{}

// StoreTransaction saves the raw and processed transaction in the datastores.
func StoreTransaction(ctx context.Context, txStore datastore.TransactionStore, rawTxStore datastore.RawTransactionStore, protocol, network string, tx *datastore.Transaction, rawTx interface{}) error {
	if bytes.Contains(tx.Data, encoding.DataPrefix()) {
		if err := rawTxStore.PutRawTransaction(ctx, protocol, network, tx.Hash, rawTx); err != nil {
			return err
		}

		if err := txStore.PutTransaction(ctx, protocol, network, tx.Hash, tx); err != nil {
			return err
		}
	}

	return nil
}
