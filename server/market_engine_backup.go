package server

import (
	"io/ioutil"
	"log"
	"os"

	"gitlab.com/around25/products/matching-engine/engine"
)

// BackupMarket saves the given snapshot of the order book as binary into the backups folder with the name of the market pair
// - It first saves into a temporary file before moving the file to the final localtion
func (mkt *marketEngine) BackupMarket(market engine.MarketBackup) error {
	file := mkt.config.config.Backup.Path + ".tmp"
	rawMarket, _ := market.ToBinary()
	ioutil.WriteFile(file, rawMarket, 0644)
	os.Rename(mkt.config.config.Backup.Path+".tmp", mkt.config.config.Backup.Path)
	return nil
}

// LoadMarketFromBackup from a backup file and update the order book with the given data
// - Also reset the Kafka partition offset to the one from the backup and replay the orders to recreate the latest state
func (mkt *marketEngine) LoadMarketFromBackup() (err error) {
	log.Printf("Loading market %s from backup file: %s", mkt.name, mkt.config.config.Backup.Path)
	file := mkt.config.config.Backup.Path
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("Backup file does not exist:", file)
		return
	}

	var market engine.MarketBackup
	log.Println("Load market from binary. Content size:", len(content))
	market.FromBinary(content)

	// load all records from the backup into the order book
	log.Println(
		"Loading market from market backup. \n",
		"Highest Bid:", market.GetHighestBid(),
		"Lowest Ask:", market.GetLowestAsk(),
		"Kafka Topic:", market.GetTopic(),
		"Kafka Partition:", market.GetPartition(),
		"Kafka Offset:", market.GetOffset(),
		"Buy Orders Found:", len(market.GetBuyOrders()),
		"Sell Orders Found:", len(market.GetSellOrders()),
	)
	mkt.LoadMarket(market)

	// mark the last message that has been processed by the engine to the one saved in the backup file
	err = mkt.consumer.SetOffset(market.Offset + 1)
	if err != nil {
		log.Fatalf("Unable to reset offset for the '%s' market on '%s' topic and partition '%d' to offset '%d'", mkt.name, market.Topic, market.Partition, market.Offset)
		return
	}

	log.Printf("Market %s loaded from backup", mkt.name)

	return nil
}

func (mkt *marketEngine) LoadMarket(market engine.MarketBackup) {
	mkt.engine.LoadMarket(market)
}
