package server

import (
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"

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

	// the next message to read is given by the market... add +1 if offset>=0 to not read the last message again
	offset := market.Offset
	if offset >= 0 {
		offset++
	}

	// load all records from the backup into the order book
	log.Info().
		Str("section", "backup").Str("action", "import").
		Str("market", mkt.name).
		Int64("offset", offset).
		Int("limit_buy_count", len(market.GetBuyOrders())).
		Int("limit_sell_count", len(market.GetSellOrders())).
		Msg("Loading market from backup")
	mkt.LoadMarket(market)
	// mark the last message that has been processed by the engine to the one saved in the backup file
	err = mkt.consumer.SetOffset(offset)
	if err != nil {
		log.Fatal().Err(err).Str("section", "backup").Str("action", "import").
			Str("market", mkt.name).
			Int64("offset", offset).Msg("Unable to reset offset")
		return
	}

	return nil
}

func (mkt *marketEngine) LoadMarket(market model.MarketBackup) {
	mkt.engine.LoadMarket(market)
}
