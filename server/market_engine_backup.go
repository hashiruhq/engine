package server

import (
	"io/ioutil"
	"log"
	"os"

	"gitlab.com/around25/products/matching-engine/model"
)

// BackupMarket saves the given snapshot of the order book as binary into the backups folder with the name of the market pair
// - It first saves into a temporary file before moving the file to the final localtion
func (mkt *marketEngine) BackupMarket(market model.MarketBackup) error {
	file := mkt.config.config.Backup.Path + ".tmp"
	rawMarket, _ := market.ToBinary()
	ioutil.WriteFile(file, rawMarket, 0644)
	os.Rename(mkt.config.config.Backup.Path+".tmp", mkt.config.config.Backup.Path)
	return nil
}

// LoadMarketFromBackup from a backup file and update the order book with the given data
// - Also reset the Kafka partition offset to the one from the backup and replay the orders to recreate the latest state
func (mkt *marketEngine) LoadMarketFromBackup() (err error) {
	file := mkt.config.config.Backup.Path
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	var market model.MarketBackup
	market.FromBinary(content)

	// load all records from the backup into the order book
	log.Printf(
		"[info] [market:%s] [init:1] Loaded market from storage with offset:%d bids:%d asks:%d\n",
		mkt.name,
		market.GetOffset(),
		len(market.GetBuyOrders()),
		len(market.GetSellOrders()),
	)
	mkt.LoadMarket(market)

	// mark the last message that has been processed by the engine to the one saved in the backup file
	err = mkt.consumer.SetOffset(market.Offset + 1)
	if err != nil {
		log.Fatalf("[fatal] [market:%s] Unable to reset offset to '%d'\n", mkt.name, market.Offset)
		return
	}

	return nil
}

func (mkt *marketEngine) LoadMarket(market model.MarketBackup) {
	mkt.engine.LoadMarket(market)
}
