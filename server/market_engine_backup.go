package server

import (
	"io/ioutil"
	"log"
	"os"
	"trading_engine/trading_engine"
)

// BackupMarket saves the given snapshot of the order book as JSON into the backups folder with the name of the market pair
// - It first saves into a temporary file before moving the file to the final localtion
func (mkt *marketEngine) BackupMarket(market trading_engine.Market) error {
	file := "/root/backups/" + mkt.name + ".json.tmp"
	rawMarket, _ := market.ToJSON()
	ioutil.WriteFile(file, rawMarket, 0644)
	os.Rename("/root/backups/"+mkt.name+".json.tmp", "/root/backups/"+mkt.name+".json")
	return nil
}

// LoadMarketFromBackup from a backup file and update the order book with the given data
// - Also reset the Kafka partition offset to the one from the backup and replay the orders to recreate the latest state
func (mkt *marketEngine) LoadMarketFromBackup() (err error) {
	log.Printf("Loading market %s from backup", mkt.name)
	file := "/root/backups/" + mkt.name + ".json"
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	var market trading_engine.Market
	market.FromJSON(content)

	// load all records from the backup into the order book
	mkt.LoadMarket(market)

	// mark the last message that has been processed by the engine to the one saved in the backup file
	err = mkt.consumer.ResetOffset(market.Topic, market.Partition, market.Offset, "")
	if err != nil {
		log.Printf("Unable to reset offset for the '%s' market on '%s' topic and partition '%d' to offset '%d'", mkt.name, market.Topic, market.Partition, market.Offset)
		return
	}

	log.Printf("Market %s loaded from backup", mkt.name)

	return nil
}

func (mkt *marketEngine) LoadMarket(market trading_engine.Market) {
	mkt.engine.LoadMarket(market)
}
