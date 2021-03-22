package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/around25/products/matching-engine/model"
	"gitlab.com/around25/products/matching-engine/utils"
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
		So(len(events), ShouldEqual, 9)
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
		So(len(events), ShouldEqual, 4)
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
		So(status.GetStatus(), ShouldEqual, model.OrderStatus_Untouched)
		So(len(events), ShouldEqual, 8)
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
		So(len(events), ShouldEqual, 9)
		So(book.GetLowestEntryPrice(), ShouldEqual, 0)
		status = events[3].GetOrderActivation()
		So(status.GetStatus(), ShouldEqual, model.OrderStatus_Pending)
	})

	Convey("Order Issue 02 - Market/Stop market ignores funds", t, func() {
		/*
			REPORTED ISSUE 24.02.2021:
			Matching engine doesn't use funds/locked _funds values during matching buy market/stop market orders.
			Use case:
			1. Create limit order for user1 to sell 2 LTC for 15 USDT
			2. Create limit order for user1 to sell 1 LTC for 100 USDT
			3. Create buy market order for user2 to buy 1 LTC (execution price = 15)
			4. Ensure that user2 has 30 USDT available on your balance
			5. Create buy market order for user2 to buy 2 LTC. As a result, user2 will receive 2 LTC and pay 115 USDT,
				 it will lead to negative balance for USDT (-85).
			--
			Correct behaviour would be to allow user2 to buy 1 LTC for 15 USDT and 0.15 LTC at price 100,
			so he will get only 1.15 LTC and pay exactly 30 USDT.
		*/
		book := NewOrderBook("ltcusd", 2, 8)
		events := make([]model.Event, 0, 5)
		order1 := model.Order{ID: 1, Price: uint64(1500), Amount: uint64(100000000), Funds: uint64(100000000), Side: sell, Type: limit, EventType: newOrder, OwnerID: 1}
		order2 := model.Order{ID: 2, Price: uint64(10000), Amount: uint64(100000000), Funds: uint64(100000000), Side: sell, Type: limit, EventType: newOrder, OwnerID: 1}
		marketOrder := model.Order{ID: 3, Price: uint64(0), Amount: uint64(200000000), Funds: uint64(3000), Side: buy, Type: market, EventType: newOrder, OwnerID: 2}
		events = events[0:0]
		book.Process(order1, &events)
		So(book.GetLowestAsk(), ShouldEqual, 1500)
		book.Process(order2, &events)
		So(book.GetLowestAsk(), ShouldEqual, 1500)
		book.Process(marketOrder, &events)

		fundsUsed := uint64(0)

		for _, e := range events {
			switch e.Type {
			case model.EventType_NewTrade:
				{
					t := e.GetTrade()
					fundsUsed += utils.Multiply(t.Price, t.Amount, 2, 8, 2)
					// log.Warn().
					// 	Int("index", i).
					// 	Uint64("price", t.Price).
					// 	Uint64("amount", t.Amount).
					// 	Uint64("bid_owner", t.BidOwnerID).
					// 	Uint64("ask_owner", t.AskOwnerID).
					// 	Msg("Trade generated")
				}
			default:
				// log.Warn().
				// 	Str("type", e.GetType().String()).
				// 	Str("market", e.Market).
				// 	Uint64("seqid", e.SeqID).
				// 	Msg("Event generated")
			}
		}
		// log.Warn().
		// 	Uint64("locked_funds", marketOrder.Funds).
		// 	Uint64("funds_used", fundsUsed).
		// 	Uint64("amount", marketOrder.Amount).
		// 	Msg("Market order after processing")

		So(fundsUsed, ShouldBeLessThanOrEqualTo, marketOrder.Funds)
		So(book.GetLowestAsk(), ShouldEqual, 10000)
		// events = events[0:0]
		// book.Process(order4, &events)
		// status := events[0].GetOrderStatus()
		// So(len(events), ShouldEqual, 6)
		// So(book.GetLowestEntryPrice(), ShouldEqual, 0)
		// status = events[3].GetOrderActivation()
		// So(status.GetStatus(), ShouldEqual, model.OrderStatus_Pending)
	})
}
