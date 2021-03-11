package engine_test

import (
	"testing"

	"gitlab.com/around25/products/matching-engine/engine"
	"gitlab.com/around25/products/matching-engine/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTradingEngineCreation(t *testing.T) {
	Convey("Given an empty trading engine", t, func() {
		tradingEngine := engine.NewTradingEngine("btcusd", 8, 8)
		events := make([]model.Event, 0, 5)
		Convey("I should be able to process a new order", func() {
			tradingEngine.Process(model.NewOrder(1, uint64(100000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			tradingEngine.Process(model.NewOrder(2, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			tradingEngine.Process(model.NewOrder(3, uint64(90000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 4)
			event := events[1]
			So(event.GetTrade().Amount, ShouldEqual, 12000000000)
			So(event.GetTrade().Price, ShouldEqual, 100000000)

			state := tradingEngine.GetOrderBook().GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 1)
		})
		Convey("I should be able to process events on the trading engine", func() {
			tradingEngine.ProcessEvent(model.NewOrder(4, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			tradingEngine.ProcessEvent(model.NewOrder(5, 0, 0, model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_CancelOrder), &events)
			tradingEngine.ProcessEvent(model.NewOrder(6, 0, 0, model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_BackupMarket), &events)
			tradingEngine.ProcessEvent(model.NewOrder(7, 0, 0, model.MarketSide_Sell, model.OrderType_Limit, 4), &events)
		})
	})
}
