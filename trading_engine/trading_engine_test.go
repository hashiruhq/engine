package trading_engine_test

import (
	"testing"

	"gitlab.com/around25/products/matching-engine/trading_engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTradingEngineCreation(t *testing.T) {
	Convey("Given an empty trading engine", t, func() {
		tradingEngine := trading_engine.NewTradingEngine()
		Convey("I should be able to process a new order", func() {
			tradingEngine.Process(trading_engine.NewOrder("TEST_1", uint64(100000000), uint64(12000000000), 1, 1, 1))
			tradingEngine.Process(trading_engine.NewOrder("TEST_2", uint64(110000000), uint64(12000000000), 2, 1, 1))
			trades := tradingEngine.Process(trading_engine.NewOrder("TEST_3", uint64(90000000), uint64(12000000000), 2, 1, 1))
			So(len(trades), ShouldEqual, 1)
			trade := trades[0]
			So(trade.Amount, ShouldEqual, 12000000000)
			So(trade.Price, ShouldEqual, 100000000)

			state := tradingEngine.GetOrderBook().GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 1)
		})
		Convey("I should be able to process events on the trading engine", func() {
			tradingEngine.ProcessEvent(trading_engine.NewOrder("EVENT_1", uint64(110000000), uint64(12000000000), 2, 1, trading_engine.EventTypeNewOrder))
			tradingEngine.ProcessEvent(trading_engine.NewOrder("TEST_3", 0, 0, 2, 1, trading_engine.EventTypeCancelOrder))
			tradingEngine.ProcessEvent(trading_engine.NewOrder("EVENT_1", 0, 0, 2, 1, trading_engine.EventTypeBackupMarket))
			tradingEngine.ProcessEvent(trading_engine.NewOrder("EVENT_1", 0, 0, 2, 1, 4))
		})
	})
}
