package datastore

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/stellar/go/historyarchive"
	"github.com/stellar/go/network"
	"github.com/stellar/go/support/errors"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/go/support/storage"
)

const (
	Pubnet  = "pubnet"
	Testnet = "testnet"
)

func CreateHistoryArchiveFromNetworkName(ctx context.Context, networkName string, userAgent string, logger *supportlog.Entry) (historyarchive.ArchiveInterface, error) {
	var historyArchiveUrls []string
	switch networkName {
	case Pubnet:
		historyArchiveUrls = network.PublicNetworkhistoryArchiveURLs
	case Testnet:
		historyArchiveUrls = network.TestNetworkhistoryArchiveURLs
	default:
		return nil, errors.Errorf("Invalid network name %s", networkName)
	}

	return historyarchive.NewArchivePool(historyArchiveUrls, historyarchive.ArchiveOptions{
		Logger: logger,
		ConnectOptions: storage.ConnectOptions{
			UserAgent: userAgent,
			Context:   ctx,
		},
	})
}

func GetLatestLedgerSequenceFromHistoryArchives(archive historyarchive.ArchiveInterface) (uint32, error) {
	has, err := archive.GetRootHAS()
	if err != nil {
		log.Error("Error getting root HAS from archives", err)
		return 0, errors.Wrap(err, "failed to retrieve the latest ledger sequence from any history archive")
	}

	return has.CurrentLedger, nil
}

func GetHistoryArchivesCheckPointFrequency() uint32 {
	// this could evolve to use other sources for checkpoint freq
	return historyarchive.DefaultCheckpointFrequency
}
