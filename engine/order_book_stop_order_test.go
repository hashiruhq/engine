package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/around25/products/matching-engine/model"
)

/**
Prod Stop orders testing

Stop=Loss
1. Create sell stop limit order: Stop=95, Limit=90
2. Create buy limit orders from other account: Limit = 96, Amount=0.2 | Limit=94, Amount=0.2
3. Create sell market order from third account: Amount=0.4 (last price is 94)
--
Expected result: Stop order triggered and executed
Actual result: Stop order is not triggered

Stop=Loss
1. Create sell stop limit order: Stop=95, Limit=90
2. Create buy limit orders from other account: Limit = 96, Amount=0.2 | Limit=94, Amount=0.2
3. Create sell market order from third account: Amount=0.2 (last price is 96)
--
Expected result: Stop order is not triggered
Actual result: Stop order is not triggered

Stop=Entry
1. Create sell stop limit order: Stop=95, Limit=90
2. Create buy limit orders from other account: Limit = 96, Amount=0.2 | Limit=94, Amount=0.2
3. Create sell market order from third account: Amount=0.4 (last price is 94)
--
Expected result: Stop order triggered and executed
Actual result: Stop order is not triggered

Stop=Entry
1. Create sell stop limit order: Stop=95, Limit=90
2. Create buy limit orders from other account: Limit = 96, Amount=0.2 | Limit=94, Amount=0.2
3. Create sell market order from third account: Amount=0.2 (last price is 96)
--
Expected result: Stop order is not triggered and executed
Actual result: Stop order triggered and executed

*/

func TestStopOrders(t *testing.T) {
	sell := model.MarketSide_Sell
	buy := model.MarketSide_Buy
	limit := model.OrderType_Limit
	market := model.OrderType_Market
	stopLoss := model.StopLoss_Loss
	stopEntry := model.StopLoss_Entry
	newOrder := model.CommandType_NewOrder

	Convey("Stop Order Test Case 1", t, func() {
		/*
			Stop=Loss
			1. Create sell stop limit order: Stop=95, Limit=90
			2. Create buy limit orders from other account: Limit = 96, Amount=0.2 | Limit=94, Amount=0.2
			3. Create sell market order from third account: Amount=0.4 (last price is 94)
			--
			Expected result: Stop order triggered and executed
			Actual result: Stop order is not triggered
		*/
		book := NewOrderBook("btcusd", 2, 8)
		events := make([]model.Event, 0, 5)
		order1 := model.Order{ID: 1, Price: uint64(9000), Amount: uint64(10000000), Funds: uint64(200000000), Side: sell, Type: limit, Stop: stopLoss, StopPrice: uint64(9500), EventType: newOrder}
		order2 := model.Order{ID: 2, Price: uint64(9600), Amount: uint64(20000000), Funds: uint64(20000000), Side: buy, Type: limit, EventType: newOrder}
		order3 := model.Order{ID: 3, Price: uint64(9400), Amount: uint64(20000000), Funds: uint64(20000000), Side: buy, Type: limit, EventType: newOrder}
		order4 := model.Order{ID: 4, Price: uint64(0), Amount: uint64(100000000), Funds: uint64(200000000), Side: sell, Type: market, EventType: newOrder}
		book.Process(order1, &events)
		So(book.GetHighestLossPrice(), ShouldEqual, 9500)
		book.Process(order2, &events)
		book.Process(order3, &events)
		events = events[0:0]
		book.Process(order4, &events)
		So(len(events), ShouldEqual, 6)
		So(book.GetHighestLossPrice(), ShouldEqual, 0)
		status := events[4].GetOrderActivation()
		So(status.GetStatus(), ShouldEqual, model.OrderStatus_Pending)
	})

	Convey("Stop Order Test Case 2", t, func() {
		/*
			Stop=Loss
			1. Create sell stop limit order: Stop=95, Limit=90
			2. Create buy limit orders from other account: Limit = 96, Amount=0.2 | Limit=94, Amount=0.2
			3. Create sell market order from third account: Amount=0.2 (last price is 96)
			--
			Expected result: Stop order is not triggered
			Actual result: Stop order is not triggered
		*/
		book := NewOrderBook("btcusd", 2, 8)
		events := make([]model.Event, 0, 5)
		order1 := model.Order{ID: 1, Price: uint64(9000), Amount: uint64(10000000), Funds: uint64(200000000), Side: sell, Type: limit, Stop: stopLoss, StopPrice: uint64(9500), EventType: newOrder}
		order2 := model.Order{ID: 2, Price: uint64(9600), Amount: uint64(20000000), Funds: uint64(20000000), Side: buy, Type: limit, EventType: newOrder}
		order3 := model.Order{ID: 3, Price: uint64(9400), Amount: uint64(20000000), Funds: uint64(20000000), Side: buy, Type: limit, EventType: newOrder}
		order4 := model.Order{ID: 4, Price: uint64(0), Amount: uint64(20000000), Funds: uint64(200000000), Side: sell, Type: market, EventType: newOrder}
		book.Process(order1, &events)
		So(book.GetHighestLossPrice(), ShouldEqual, 9500)
		book.Process(order2, &events)
		book.Process(order3, &events)
		events = events[0:0]
		book.Process(order4, &events)
		So(len(events), ShouldEqual, 3)
		So(book.GetHighestLossPrice(), ShouldEqual, 9500)
	})

	Convey("Stop Order Test Case 3", t, func() {
		/*
			Stop=Entry
			1. Create sell stop limit order: Stop=95, Limit=90
			2. Create buy limit orders from other account: Limit = 96, Amount=0.2 | Limit=94, Amount=0.2
			3. Create sell market order from third account: Amount=0.4 (last price is 94)
			--
			Expected result: Stop order triggered and executed
			Actual result: Stop order is not triggered
		*/
		book := NewOrderBook("btcusd", 2, 8)
		events := make([]model.Event, 0, 5)
		order1 := model.Order{ID: 1, Price: uint64(9000), Amount: uint64(10000000), Funds: uint64(200000000), Side: sell, Type: limit, Stop: stopEntry, StopPrice: uint64(9500), EventType: newOrder}
		order2 := model.Order{ID: 2, Price: uint64(9600), Amount: uint64(20000000), Funds: uint64(1920), Side: buy, Type: limit, EventType: newOrder}
		order3 := model.Order{ID: 3, Price: uint64(9400), Amount: uint64(20000000), Funds: uint64(1880), Side: buy, Type: limit, EventType: newOrder}
		order4 := model.Order{ID: 4, Price: uint64(0), Amount: uint64(40000000), Funds: uint64(200000000), Side: sell, Type: market, EventType: newOrder}
		book.Process(order1, &events)
		So(book.GetLowestEntryPrice(), ShouldEqual, 9500)
		book.Process(order2, &events)
		book.Process(order3, &events)
		events = events[0:0]
		book.Process(order4, &events)
		status := events[0].GetOrderStatus()
		So(len(events), ShouldEqual, 6)
		So(book.GetLowestEntryPrice(), ShouldEqual, 0)
		status = events[4].GetOrderActivation()
		So(status.GetStatus(), ShouldEqual, model.OrderStatus_Pending)
	})

	Convey("Stop Order Test Case 4", t, func() {
		/*
			Stop=Entry
			1. Create sell stop limit order: Stop=95, Limit=90
			2. Create buy limit orders from other account: Limit = 96, Amount=0.2 | Limit=94, Amount=0.2
			3. Create sell market order from third account: Amount=0.2 (last price is 96)
			--
			NOT RIGHT: Expected result: Stop order is not triggered and executed
			RIGHT: Actual result: Stop order triggered and executed
		*/
		book := NewOrderBook("btcusd", 2, 8)
		events := make([]model.Event, 0, 5)
		order1 := model.Order{ID: 1, Price: uint64(9000), Amount: uint64(10000000), Funds: uint64(200000000), Side: sell, Type: limit, Stop: stopEntry, StopPrice: uint64(9500), EventType: newOrder}
		order2 := model.Order{ID: 2, Price: uint64(9600), Amount: uint64(20000000), Funds: uint64(1920), Side: buy, Type: limit, EventType: newOrder}
		order3 := model.Order{ID: 3, Price: uint64(9400), Amount: uint64(20000000), Funds: uint64(1880), Side: buy, Type: limit, EventType: newOrder}
		order4 := model.Order{ID: 4, Price: uint64(0), Amount: uint64(20000000), Funds: uint64(200000000), Side: sell, Type: market, EventType: newOrder}
		events = events[0:0]
		book.Process(order1, &events)
		So(book.GetLowestEntryPrice(), ShouldEqual, 9500)
		book.Process(order2, &events)
		book.Process(order3, &events)
		events = events[0:0]
		book.Process(order4, &events)
		status := events[0].GetOrderStatus()
		So(len(events), ShouldEqual, 6)
		So(book.GetLowestEntryPrice(), ShouldEqual, 0)
		status = events[3].GetOrderActivation()
		So(status.GetStatus(), ShouldEqual, model.OrderStatus_Pending)
	})
}
