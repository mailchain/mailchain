package processor

import (
	"context"
	"fmt"

	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	"github.com/mailchain/mailchain/cmd/indexer/internal/clients"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
)

type Sequential struct {
	protocol string
	network  string

	syncStore      datastore.SyncStore
	blockProcessor actions.Block
	blockClient    clients.BlockByNumber
}

func NewSequential(protocol, network string, store datastore.SyncStore, proc actions.Block, client clients.BlockByNumber) *Sequential {
	return &Sequential{
		protocol:       protocol,
		network:        network,
		syncStore:      store,
		blockProcessor: proc,
		blockClient:    client,
	}
}

func (s *Sequential) NextBlock(ctx context.Context) error {
	blkNo, err := s.syncStore.GetBlockNumber(ctx, s.protocol, s.network)
	if err != nil {
		return err
	}

	blk, err := s.blockClient.Get(ctx, blkNo)
	if err != nil {
		return err
	}

	if err := s.blockProcessor.Run(ctx, s.protocol, s.network, blk); err != nil {
		return err
	}

	nextBlockNo := blkNo + 1

	fmt.Println("completed block number: ", blkNo)

	if err := s.syncStore.PutBlockNumber(ctx, s.protocol, s.network, nextBlockNo); err != nil {
		return err
	}

	return nil
}
