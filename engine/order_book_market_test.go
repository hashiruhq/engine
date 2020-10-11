package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/around25/products/matching-engine/model"
	"gitlab.com/around25/products/matching-engine/utils"
)

func TestBuyMarketOrderIssue(t *testing.T) {
	events := make([]model.Event, 0, 5)

	Convey("Check all order book fields", t, func() {
		orderBook := NewOrderBook("ethusdt", 5, 1)
		So(orderBook.GetMarketID(), ShouldEqual, "ethusdt")
		So(orderBook.GetPricePrecision(), ShouldEqual, 5)
		So(orderBook.GetVolumePrecision(), ShouldEqual, 1)
		So(orderBook.GetLastEventSeqID(), ShouldEqual, 0)
		So(orderBook.GetLastTradeSeqID(), ShouldEqual, 0)
		So(orderBook.GetHighestLossPrice(), ShouldEqual, 0)
		So(orderBook.GetLowestEntryPrice(), ShouldEqual, 0)
	})

	Convey("Given an empty order book", t, func() {

		Convey("Issue 1: Use case 1/2: I should be able to add a market buy order in an existing market", func() {
			orderBook := NewOrderBook("prdxusdt", 5, 1)

			events = events[0:0]
			order1 := model.NewOrder(90664015, uint64(12255), uint64(11045), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder)
			orderBook.Process(order1, &events)

			order2 := model.NewOrder(90664015, uint64(12259), uint64(10605), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder)
			events = events[0:0]
			orderBook.Process(order2, &events)

			events = events[0:0]
			lockedFunds := uint64(19900000)
			marketOrder := model.Order{ID: 90664126, Amount: 16704, Funds: lockedFunds, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}
			orderBook.Process(marketOrder, &events)

			So(len(events), ShouldEqual, 4)
			So(events[0].GetOrderStatus().GetStatus(), ShouldEqual, model.OrderStatus_Untouched)
			trade1 := events[1].GetTrade()
			trade2 := events[2].GetTrade()
			So(events[3].GetOrderStatus().GetStatus(), ShouldEqual, model.OrderStatus_Cancelled)

			usedFunds := trade1.Price*trade1.Amount/10 + trade2.Price*trade2.Amount/10

			So(lockedFunds, ShouldBeGreaterThanOrEqualTo, usedFunds)
			So(trade1.Price, ShouldEqual, 12255)
			So(trade1.Amount, ShouldEqual, 11045)
			So(trade2.Price, ShouldEqual, 12259)
			So(trade2.Amount, ShouldEqual, 5191)
		})

		Convey("Issue 1: Use case 2/2: I should be able to add a market buy order in an existing market", func() {
			orderBook := NewOrderBook("ethusdt", 2, 5)
			lockedFunds := uint64(66058)

			events = events[0:0]
			order1 := model.NewOrder(92454370, uint64(19034), uint64(16434), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder)
			orderBook.Process(order1, &events)

			order2 := model.NewOrder(92454372, uint64(19035), uint64(21610), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder)
			events = events[0:0]
			orderBook.Process(order2, &events)

			order3 := model.NewOrder(92454373, uint64(19036), uint64(400000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder)
			events = events[0:0]
			orderBook.Process(order3, &events)

			marketOrder := model.Order{ID: 92454464, Price: 18999, Amount: 347000, StopPrice: 18999, Funds: lockedFunds, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}
			events = events[0:0]
			orderBook.Process(marketOrder, &events)

			So(len(events), ShouldEqual, 5)
			So(events[0].GetOrderStatus().GetStatus(), ShouldEqual, model.OrderStatus_Untouched)
			t1 := events[1].GetTrade()
			t2 := events[2].GetTrade()
			t3 := events[3].GetTrade()
			usedFunds := utils.Multiply(t1.Price, t1.Amount, 2, 5, 2) +
				utils.Multiply(t2.Price, t2.Amount, 2, 5, 2) +
				utils.Multiply(t3.Price, t3.Amount, 2, 5, 2)
			So(events[4].GetOrderStatus().GetStatus(), ShouldEqual, model.OrderStatus_Cancelled)
			So(events[4].GetOrderStatus().GetFunds(), ShouldEqual, lockedFunds-usedFunds)

			So(lockedFunds, ShouldBeGreaterThanOrEqualTo, usedFunds)
			So(t1.Price, ShouldEqual, 19034)
			So(t1.Amount, ShouldEqual, 16434)
			So(t2.Price, ShouldEqual, 19035)
			So(t2.Amount, ShouldEqual, 21610)
			So(t3.Price, ShouldEqual, 19036)
			So(t3.Amount, ShouldEqual, 308956)
		})

		// Test price rounding for market order
		// Market Buy: Compare engine logic for market orders with other dev engines
		Convey("Market Buy: Compare engine logic for market orders with other dev engines", func() {
			orderBook := NewOrderBook("btcusdt", 2, 8)
			lockedFunds := uint64(500000)

			events = events[0:0]
			order1 := model.NewOrder(1, uint64(100000), uint64(5000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder)
			orderBook.Process(order1, &events)

			events = events[0:0]
			marketOrder := model.Order{ID: 2, Price: 0, Amount: 4997355, Funds: lockedFunds, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}
			orderBook.Process(marketOrder, &events)

			So(len(events), ShouldEqual, 3)
			So(events[0].GetOrderStatus().GetStatus(), ShouldEqual, model.OrderStatus_Untouched)
			trade1 := events[1].GetTrade()
			So(events[2].GetOrderStatus().GetStatus(), ShouldEqual, model.OrderStatus_Cancelled)
			So(events[2].GetOrderStatus().GetFunds(), ShouldEqual, 500000-4998)

			So(lockedFunds, ShouldBeGreaterThanOrEqualTo, 4998)
			So(trade1.Price, ShouldEqual, 100000)
			So(trade1.Amount, ShouldEqual, 4997355)
		})

		Convey("I should be able to add a market buy order in an existing market 2", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(32, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 33, Amount: 10000000000, Funds: 1100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 3)
			So(events[1].GetTrade().Price, ShouldEqual, 110000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 1000000000)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 34, Amount: 15000000000, Funds: 23456000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 3)
			So(events[1].GetTrade().Price, ShouldEqual, 110000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 3", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(35, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(36, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 37, Amount: 10000000000, Funds: 1100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 3)
			So(events[1].GetTrade().Price, ShouldEqual, 110000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 1000000000)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 38, Amount: 15000000000, Funds: 23456000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 4)
			So(events[1].GetTrade().Price, ShouldEqual, 110000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 4", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(39, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(40, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 41, Amount: 12000000000, Funds: 51100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 3)
			So(events[1].GetTrade().Price, ShouldEqual, 110000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market buy order in an empty market", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			events = events[0:0]
			orderBook.Process(model.Order{ID: 59, Amount: 11000000000, Funds: 1000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 2)
		})

		Convey("I should be able to add a market buy order after another pending market order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(61, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.Order{ID: 62, Amount: 14000000000, Funds: 20000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			events = events[0:0]
			orderBook.Process(model.Order{ID: 63, Amount: 11000000000, Funds: 1000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 2)
		})

		Convey("I should be able to add a limit sell order after another pending market buy order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 67, Amount: 11000000000, Funds: 13200000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(68, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetOrderStatus().Status, ShouldEqual, model.OrderStatus_Untouched)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 71, Amount: 11000000000, Funds: 13200000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			orderBook.Process(model.Order{ID: 72, Amount: 11000000000, Funds: 13200000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(73, uint64(120000000), uint64(24000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetOrderStatus().Status, ShouldEqual, model.OrderStatus_Untouched)
		})

		Convey("I should be able to cancel an invalid order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			order := model.NewOrder(96, uint64(110000000), uint64(800000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder)
			events = events[0:0]
			orderBook.Cancel(order, &events)
			So(len(events), ShouldEqual, 1)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)

			order = model.Order{ID: 97, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}
			events = events[0:0]
			orderBook.Cancel(order, &events)
			So(len(events), ShouldEqual, 0)

			order = model.Order{ID: 98, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}
			events = events[0:0]
			orderBook.Cancel(order, &events)
			So(len(events), ShouldEqual, 0)

			order = model.Order{ID: 99, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: 4, Side: model.MarketSide_Buy}
			events = events[0:0]
			orderBook.Cancel(order, &events)
			So(len(events), ShouldEqual, 0)
		})

	})
}
