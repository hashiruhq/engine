package trading_engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTradingEngineCreation(t *testing.T) {
	Convey("Given an empty trading engine", t, func() {
		tradingEngine := NewTradingEngine()
		Convey("I should be able to process a new order", func() {
			tradingEngine.Process(NewOrder("TEST_1", uint64(100000000), uint64(12000000000), 1, 1))
			tradingEngine.Process(NewOrder("TEST_2", uint64(110000000), uint64(12000000000), 2, 1))
			trades := tradingEngine.Process(NewOrder("TEST_3", uint64(90000000), uint64(12000000000), 2, 1))
			So(len(trades), ShouldEqual, 1)
			trade := trades[0]
			So(trade.Amount, ShouldEqual, 12000000000)
			So(trade.Price, ShouldEqual, 100000000)
			So(tradingEngine.GetOrderBook().GetMarket().Len(), ShouldEqual, 1)
		})
	})
}
