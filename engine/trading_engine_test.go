package engine_test

import (
	"testing"

	"gitlab.com/around25/products/matching-engine/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTradingEngineCreation(t *testing.T) {
	Convey("Given an empty trading engine", t, func() {
		tradingEngine := engine.NewTradingEngine()
		Convey("I should be able to process a new order", func() {
			tradingEngine.Process(engine.NewOrder("TEST_1", uint64(100000000), uint64(12000000000), engine.MarketSide_Buy, engine.OrderType_Limit, engine.CommandType_NewOrder))
			tradingEngine.Process(engine.NewOrder("TEST_2", uint64(110000000), uint64(12000000000), engine.MarketSide_Sell, engine.OrderType_Limit, engine.CommandType_NewOrder))
			trades := tradingEngine.Process(engine.NewOrder("TEST_3", uint64(90000000), uint64(12000000000), engine.MarketSide_Sell, engine.OrderType_Limit, engine.CommandType_NewOrder))
			So(len(trades), ShouldEqual, 1)
			trade := trades[0]
			So(trade.Amount, ShouldEqual, 12000000000)
			So(trade.Price, ShouldEqual, 100000000)

			state := tradingEngine.GetOrderBook().GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 1)
		})
		Convey("I should be able to process events on the trading engine", func() {
			tradingEngine.ProcessEvent(engine.NewOrder("EVENT_1", uint64(110000000), uint64(12000000000), engine.MarketSide_Sell, engine.OrderType_Limit, engine.CommandType_NewOrder))
			tradingEngine.ProcessEvent(engine.NewOrder("TEST_3", 0, 0, engine.MarketSide_Sell, engine.OrderType_Limit, engine.CommandType_CancelOrder))
			tradingEngine.ProcessEvent(engine.NewOrder("EVENT_1", 0, 0, engine.MarketSide_Sell, engine.OrderType_Limit, engine.CommandType_BackupMarket))
			tradingEngine.ProcessEvent(engine.NewOrder("EVENT_1", 0, 0, engine.MarketSide_Sell, engine.OrderType_Limit, 4))
		})
	})
}
