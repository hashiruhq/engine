package engine_test

import (
	"testing"

	"gitlab.com/around25/products/matching-engine/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTradingEngineCreation(t *testing.T) {
	Convey("Given an empty trading engine", t, func() {
		tradingEngine := engine.NewTradingEngine("btcusd", 8, 8)
		trades := make([]engine.Trade, 0, 5)
		Convey("I should be able to process a new order", func() {
			tradingEngine.Process(engine.NewOrder(1, uint64(100000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &trades)
			tradingEngine.Process(engine.NewOrder(2, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &trades)
			trades = trades[0:0]
			tradingEngine.Process(engine.NewOrder(3, uint64(90000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			trade := trades[0]
			So(trade.Amount, ShouldEqual, 12000000000)
			So(trade.Price, ShouldEqual, 100000000)

			state := tradingEngine.GetOrderBook().GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 1)
		})
		Convey("I should be able to process events on the trading engine", func() {
			tradingEngine.ProcessEvent(engine.NewOrder(4, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &trades)
			tradingEngine.ProcessEvent(engine.NewOrder(5, 0, 0, model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_CancelOrder), &trades)
			tradingEngine.ProcessEvent(engine.NewOrder(6, 0, 0, model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_BackupMarket), &trades)
			tradingEngine.ProcessEvent(engine.NewOrder(7, 0, 0, model.MarketSide_Sell, model.OrderType_Limit, 4), &trades)
		})
	})
}
