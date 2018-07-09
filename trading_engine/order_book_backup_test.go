package trading_engine_test

import (
	"testing"
	"trading_engine/trading_engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOrderBookBackupWithMarket(t *testing.T) {
	Convey("Should be able to load an order book from a market record", t, func() {
		engine := trading_engine.NewTradingEngine()
		book := engine.GetOrderBook()
		book.Process(trading_engine.NewOrder("TEST_1", uint64(100000000), uint64(12000000000), 1, 1, 1))
		book.Process(trading_engine.NewOrder("TEST_2", uint64(110000000), uint64(12000000000), 2, 1, 1))

		market := engine.BackupMarket()
		So(market.LowestAsk, ShouldEqual, 110000000)
		So(market.HighestBid, ShouldEqual, 100000000)

		engine2 := trading_engine.NewTradingEngine()
		engine2.LoadMarket(market)
		book2 := engine2.GetOrderBook()

		So(book.GetHighestBid(), ShouldEqual, book2.GetHighestBid())
		So(book.GetLowestAsk(), ShouldEqual, book2.GetLowestAsk())
	})
}
