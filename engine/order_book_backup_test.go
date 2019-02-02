package engine_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOrderBookBackupWithMarket(t *testing.T) {
	Convey("Should be able to load an order book from a market record", t, func() {
		// ngin := engine.NewTradingEngine()
		// book := ngin.GetOrderBook()
		// book.Process(engine.NewOrder("TEST_1", uint64(100000000), uint64(12000000000), 1, 1, 1))
		// book.Process(engine.NewOrder("TEST_2", uint64(110000000), uint64(12000000000), 2, 1, 1))

		// market := ngin.BackupMarket()
		// So(market.LowestAsk, ShouldEqual, 110000000)
		// So(market.HighestBid, ShouldEqual, 100000000)

		// ngin2 := engine.NewTradingEngine()
		// ngin2.LoadMarket(market)
		// book2 := ngin.GetOrderBook()

		// So(book.GetHighestBid(), ShouldEqual, book2.GetHighestBid())
		// So(book.GetLowestAsk(), ShouldEqual, book2.GetLowestAsk())
	})
}
