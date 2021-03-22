package engine

import (
	"testing"

	"gitlab.com/around25/products/matching-engine/model"

	. "github.com/smartystreets/goconvey/convey"
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

func TestOrderBookProcessing(t *testing.T) {
	book := NewOrderBook("btcusd", 8, 8)
	events := make([]model.Event, 0, 5)

	Convey("Check all order book fields", t, func() {
		orderBook := NewOrderBook("ethbtc", 5, 8)
		So(orderBook.GetMarketID(), ShouldEqual, "ethbtc")
		So(orderBook.GetPricePrecision(), ShouldEqual, 5)
		So(orderBook.GetVolumePrecision(), ShouldEqual, 8)
		So(orderBook.GetLastEventSeqID(), ShouldEqual, 0)
		So(orderBook.GetLastTradeSeqID(), ShouldEqual, 0)
		So(orderBook.GetHighestLossPrice(), ShouldEqual, 0)
		So(orderBook.GetLowestEntryPrice(), ShouldEqual, 0)
	})

	Convey("Given an empty order book", t, func() {
		events = events[0:0]
		//	BUY          SELL
		//	1.0 120      -
		Convey("Add a first buy order", func() {
			book.Process(model.NewOrder(1, uint64(100000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
		})
		// BUY          SELL
		// 1.0 120      1.1 120
		Convey("Add a first sell order", func() {
			book.Process(model.NewOrder(2, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
		})
		// BUY          SELL
		// -            1.1 120
		Convey("Add a matching sell order", func() {
			book.Process(model.NewOrder(3, uint64(90000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 4)
			So(events[1].GetTrade().Amount, ShouldEqual, 12000000000)
			So(events[1].GetTrade().Price, ShouldEqual, 100000000)

			state := book.GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 1)
		})
		// ORDER: BUY 1.1 20
		// BUY          SELL
		// -            1.1 100
		Convey("Add a sell order with the same price", func() {
			book.Process(model.NewOrder(4, uint64(110000000), uint64(2000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 4)
			So(events[1].GetTrade().Amount, ShouldEqual, 2000000000)
			So(events[1].GetTrade().Price, ShouldEqual, 110000000)

			So(book.GetLowestAsk(), ShouldEqual, 110000000)
		})
		// ORDER: BUY 1.11 120
		// BUY          SELL
		// 1.11 20      -
		Convey("Add a buy order with a larger amount than the available sell", func() {
			book.Process(model.NewOrder(5, uint64(111000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 4)
			So(events[1].GetTrade().Amount, ShouldEqual, 10000000000)
			So(events[1].GetTrade().Price, ShouldEqual, 110000000)
			So(book.GetHighestBid(), ShouldEqual, 111000000)

			state := book.GetMarket()
			So(state[0].Len(), ShouldEqual, 1)
			So(state[1].Len(), ShouldEqual, 0)
		})
		// ORDER: BUY 1.20 120
		// BUY          SELL
		// 1.20 120      -
		// 1.11 20
		Convey("Add two another buy orders with a higher price", func() {
			book.Process(model.NewOrder(6, uint64(120000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			events = events[0:0]
			book.Process(model.NewOrder(7, uint64(120000000), uint64(1000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
		})
		// ORDER: SELL 1.30 10
		// BUY          SELL
		// 1.20 130     1.30 10
		// 1.11 20
		Convey("Add two sell orders at the same price without matching", func() {
			book.Process(model.NewOrder(7, uint64(130000000), uint64(1000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			events = events[0:0]
			book.Process(model.NewOrder(8, uint64(130000000), uint64(1000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
		})
		// ORDER: SELL 1.40 20
		// BUY          SELL
		// 1.20 130     -
		// 1.11 20
		Convey("Add a buy order that clears the sell side of the order book", func() {
			book.Process(model.NewOrder(9, uint64(140000000), uint64(2000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 6)
		})
		// ORDER: SELL 1.00 140
		// BUY          SELL
		// -            -
		Convey("Add a sell order that clears the buy side of the order book", func() {
			book.Process(model.NewOrder(10, uint64(100000000), uint64(15000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 8)

			state := book.GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 0)
		})
		Convey("Lowest ask price should be updated", func() {
			book.Process(model.NewOrder(11, uint64(110000000), uint64(1000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			book.Process(model.NewOrder(12, uint64(130000000), uint64(1000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			book.Process(model.NewOrder(13, uint64(140000000), uint64(1300000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 6)
			So(book.GetLowestAsk(), ShouldEqual, uint64(130000000))
		})

		Convey("Highest ask should be 0 when there are no more buy orders", func() {
			book.Process(model.NewOrder(14, uint64(140000000), uint64(1000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			book.Process(model.NewOrder(15, uint64(110000000), uint64(800000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 4)
			So(book.GetHighestBid(), ShouldEqual, 0)
		})

		Convey("Check if adding multiple orders works as expected", func() {
			/*
			   BUY        SELL
			   1.21 70    1.13 140 -> 1.11 10     1.13 70
			   1.11 10    1.16 70
			*/

			//buy
			book.Process(model.NewOrder(16, uint64(111000000), uint64(100000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			book.Process(model.NewOrder(17, uint64(121000000), uint64(700000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			//sell
			book.Process(model.NewOrder(18, uint64(116000000), uint64(700000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			book.Process(model.NewOrder(19, uint64(113000000), uint64(700000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			book.Process(model.NewOrder(20, uint64(113000000), uint64(700000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			backup := book.Backup()

			So(backup.HighestBid, ShouldBeLessThan, backup.LowestAsk)
		})

		Convey("Check if the lowest ask is moved after a completed limit order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(20, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(21, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(22, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 4)
			So(events[1].GetTrade().Price, ShouldEqual, 110000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 12000000000)
		})

		Convey("Check that the price of generated trades is within buy/sell limit bounds", func() {
			// cleanup
			events = events[0:0]
			// in a new market
			orderBook := NewOrderBook("ethbtc", 8, 5)
			// add a 947@200000 limit sell order
			sell01 := model.Order{ID: 20, Price: uint64(200000), Amount: uint64(947), Funds: uint64(947), Side: model.MarketSide_Sell, Type: model.OrderType_Limit, EventType: model.CommandType_NewOrder}
			orderBook.Process(sell01, &events)
			So(orderBook.GetLowestAsk(), ShouldEqual, 200000)
			So(orderBook.GetHighestBid(), ShouldEqual, 0)
			// add another 947@113000 limit sell order
			sell02 := model.Order{ID: 21, Price: uint64(113000), Amount: uint64(947), Funds: uint64(947), Side: model.MarketSide_Sell, Type: model.OrderType_Limit, EventType: model.CommandType_NewOrder}
			orderBook.Process(sell02, &events)
			// check levels
			So(orderBook.GetLowestAsk(), ShouldEqual, 113000)
			So(orderBook.GetHighestBid(), ShouldEqual, 0)
			So(orderBook.GetHighestLossPrice(), ShouldEqual, 0)
			So(orderBook.GetLowestEntryPrice(), ShouldEqual, 0)

			// stop order [100@66000 after 113000]
			events = events[0:0]
			stop01 := model.Order{
				ID:        22,
				Price:     uint64(66000),
				Amount:    uint64(100),
				Funds:     uint64(154000),
				Side:      model.MarketSide_Buy,
				Type:      model.OrderType_Limit,
				Stop:      model.StopLoss_Loss,
				StopPrice: uint64(113000),
				EventType: model.CommandType_NewOrder,
			}
			orderBook.Process(stop01, &events)

			So(orderBook.GetHighestLossPrice(), ShouldEqual, 113000)
			So(orderBook.GetLowestEntryPrice(), ShouldEqual, 0)

			// add 234@133000 limit buy order
			// => trade BUY 234@113000
			// => order partially filled SELL 713@113000
			// => stop order activated
			//    -> add limit order BUY 100@66000, untouched
			events = events[0:0]
			buy01 := model.Order{ID: 23, Price: uint64(113000), Amount: uint64(234), Funds: uint64(154000), Side: model.MarketSide_Buy, Type: model.OrderType_Limit, EventType: model.CommandType_NewOrder}
			orderBook.Process(buy01, &events)
			So(len(events), ShouldEqual, 6)
			So(orderBook.GetLowestAsk(), ShouldEqual, 113000)
			So(orderBook.GetHighestBid(), ShouldEqual, 66000)
			So(orderBook.GetHighestLossPrice(), ShouldEqual, 0)
			So(orderBook.GetLowestEntryPrice(), ShouldEqual, 0)
			So(events[0].GetOrderStatus().Price, ShouldEqual, 113000)
			So(events[0].GetOrderStatus().Amount, ShouldEqual, 234)
			So(events[1].GetTrade().Price, ShouldEqual, 113000)
			So(events[1].GetTrade().Amount, ShouldEqual, 234)
			So(events[4].GetOrderActivation().Price, ShouldEqual, 66000)
			So(events[4].GetOrderActivation().Amount, ShouldEqual, 100)
			So(events[5].GetOrderStatus().Price, ShouldEqual, 66000)
			So(events[5].GetOrderStatus().Amount, ShouldEqual, 100)

			// cancel activated stop order resulting in order book:
			// SELL 947@200000
			// SELL 713@113000
			// ---------------
			events = events[0:0]
			cStop01 := model.Order{
				ID:        22,
				Price:     uint64(66000),
				Amount:    uint64(100),
				Funds:     uint64(154000),
				Side:      model.MarketSide_Buy,
				Type:      model.OrderType_Limit,
				Stop:      model.StopLoss_Loss,
				StopPrice: uint64(113000),
				EventType: model.CommandType_CancelOrder,
			}
			orderBook.Process(cStop01, &events)
			So(len(events), ShouldEqual, 1)
			So(orderBook.GetLowestAsk(), ShouldEqual, 113000)
			So(orderBook.GetHighestBid(), ShouldEqual, 0)
			So(orderBook.GetHighestLossPrice(), ShouldEqual, 0)
			So(orderBook.GetLowestEntryPrice(), ShouldEqual, 0)

			// and another stop order BUY 13@66000 after113000
			events = events[0:0]
			stop02 := model.Order{
				ID:        24,
				Price:     uint64(66000),
				Amount:    uint64(13),
				Funds:     uint64(154000),
				Side:      model.MarketSide_Buy,
				Type:      model.OrderType_Limit,
				Stop:      model.StopLoss_Loss,
				StopPrice: uint64(113000),
				EventType: model.CommandType_NewOrder,
			}
			orderBook.Process(stop02, &events)
			So(len(events), ShouldEqual, 1)
			So(orderBook.GetHighestLossPrice(), ShouldEqual, 113000)

			// cancel the last stop buy order
			events = events[0:0]
			orderBook.Process(model.Order{
				ID:        24,
				Price:     uint64(66000),
				Amount:    uint64(13),
				Side:      model.MarketSide_Buy,
				Stop:      model.StopLoss_Loss,
				StopPrice: uint64(113000),
				Type:      model.OrderType_Limit,
				EventType: model.CommandType_CancelOrder,
			}, &events)
			So(len(events), ShouldEqual, 1)
			So(orderBook.GetHighestLossPrice(), ShouldEqual, 0)
			So(orderBook.GetLowestEntryPrice(), ShouldEqual, 0)

			// and another stop order SELL 13@66000 after113000
			events = events[0:0]
			stop03 := model.Order{
				ID:        25,
				Price:     uint64(66000),
				Amount:    uint64(13),
				Funds:     uint64(154000),
				Side:      model.MarketSide_Sell,
				Type:      model.OrderType_Limit,
				Stop:      model.StopLoss_Entry,
				StopPrice: uint64(113000),
				EventType: model.CommandType_NewOrder,
			}
			orderBook.Process(stop03, &events)
			So(len(events), ShouldEqual, 1)
			So(orderBook.GetLowestEntryPrice(), ShouldEqual, 113000)

			events = events[0:0]
			stop04 := model.Order{
				ID:        251,
				Side:      model.MarketSide_Sell,
				Type:      model.OrderType_Limit,
				Stop:      model.StopLoss_Entry,
				StopPrice: uint64(114000),
				EventType: model.CommandType_NewOrder,
			}
			orderBook.Process(stop04, &events)
			// cancel the last stop buy order
			events = events[0:0]
			orderBook.Process(model.Order{
				ID:        25,
				Side:      model.MarketSide_Sell,
				Stop:      model.StopLoss_Entry,
				StopPrice: uint64(113000),
				Type:      model.OrderType_Limit,
				EventType: model.CommandType_CancelOrder,
			}, &events)
			So(len(events), ShouldEqual, 1)
			So(orderBook.GetHighestLossPrice(), ShouldEqual, 0)
			So(orderBook.GetLowestEntryPrice(), ShouldEqual, 114000)
			events = events[0:0]
			orderBook.Process(model.Order{
				ID:        251,
				Side:      model.MarketSide_Sell,
				Stop:      model.StopLoss_Entry,
				StopPrice: uint64(114000),
				Type:      model.OrderType_Limit,
				EventType: model.CommandType_CancelOrder,
			}, &events)
			So(len(events), ShouldEqual, 1)
			So(orderBook.GetHighestLossPrice(), ShouldEqual, 0)
			So(orderBook.GetLowestEntryPrice(), ShouldEqual, 0)

			// add another buy order
			events = events[0:0]
			buy02 := model.Order{ID: 26, Price: uint64(113000), Amount: uint64(713), Funds: uint64(154000), Side: model.MarketSide_Buy, Type: model.OrderType_Limit, EventType: model.CommandType_NewOrder}
			orderBook.Process(buy02, &events)
			So(len(events), ShouldEqual, 4)
			So(orderBook.GetLowestAsk(), ShouldEqual, 200000)
			So(orderBook.GetHighestBid(), ShouldEqual, 0)
			So(events[0].GetOrderStatus().Price, ShouldEqual, 113000)
			So(events[0].GetOrderStatus().Amount, ShouldEqual, 713)
			So(events[1].GetTrade().Price, ShouldEqual, 113000)
			So(events[1].GetTrade().Amount, ShouldEqual, 713)
			// add another buy order
			events = events[0:0]
			buy03 := model.Order{ID: 27, Price: uint64(66000), Amount: uint64(234), Funds: uint64(154000), Side: model.MarketSide_Buy, Type: model.OrderType_Limit, EventType: model.CommandType_NewOrder}
			orderBook.Process(buy03, &events)
			So(len(events), ShouldEqual, 1)
			So(orderBook.GetLowestAsk(), ShouldEqual, 200000)
			So(orderBook.GetHighestBid(), ShouldEqual, 66000)
			So(events[0].GetOrderStatus().Price, ShouldEqual, 66000)
			So(events[0].GetOrderStatus().Amount, ShouldEqual, 234)

			// cancel the last buy order
			events = events[0:0]
			orderBook.Process(model.Order{ID: 27, Price: uint64(66000), Side: model.MarketSide_Buy, Type: model.OrderType_Limit, EventType: model.CommandType_CancelOrder}, &events)

			// add the last order again
			events = events[0:0]
			buy04 := model.Order{ID: 28, Price: uint64(66000), Amount: uint64(234), Funds: uint64(154000), Side: model.MarketSide_Buy, Type: model.OrderType_Limit, EventType: model.CommandType_NewOrder}
			orderBook.Process(buy04, &events)
			So(len(events), ShouldEqual, 1)
			So(orderBook.GetLowestAsk(), ShouldEqual, 200000)
			So(orderBook.GetHighestBid(), ShouldEqual, 66000)
			So(events[0].GetOrderStatus().Price, ShouldEqual, 66000)
			So(events[0].GetOrderStatus().Amount, ShouldEqual, 234)
		})

		Convey("Check if the highest bid is moved after a completed limit order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(23, uint64(120000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(24, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(25, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 4)
			So(events[1].GetTrade().Price, ShouldEqual, 120000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 12000000000)
		})

		Convey("Check if events are returned if the pricepoint is not complete", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(26, uint64(120000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(27, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(28, uint64(110000000), uint64(11000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 4)
			So(events[1].GetTrade().Price, ShouldEqual, 120000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to cancel a sell order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			order := model.NewOrder(92, uint64(110000000), uint64(800000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder)
			orderBook.Process(order, &events)
			events = events[0:0]
			orderBook.Cancel(order, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetOrderStatus().ID, ShouldEqual, order.ID)
			So(events[0].GetOrderStatus().Status, ShouldEqual, model.OrderStatus_Cancelled)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
		})

		Convey("I should be able to cancel a buy order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			order := model.NewOrder(93, uint64(110000000), uint64(800000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder)
			orderBook.Process(order, &events)
			events = events[0:0]
			orderBook.Cancel(order, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetOrderStatus().ID, ShouldEqual, order.ID)
			So(events[0].GetOrderStatus().Status, ShouldEqual, model.OrderStatus_Cancelled)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
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

		/**
		Bug report 22.03.2021: Duplicate statuses generated for these orders:
		INSERT INTO public.orders
		(id, owner_id, market_id, type, side, status, created_at, updated_at, stop, price, amount, locked_funds, filled_amount, funds, fee_amount, stop_price, used_funds, oto_type, parent_order_id, tp_price, tp_amount, tp_filled_amount, tp_status, sl_price, sl_amount, sl_filled_amount, sl_status, tp_rel_price, sl_rel_price, is_mm, init_order_id, tp_order_id, sl_order_id, account_group, filled_opposite_amount, ref_id) VALUES
		(1, 5, 'btcusdt', 'limit', 'buy', 'filled', '2021-03-17 12:44:33.524140', '2021-03-17 12:44:37.198435', 'none', 50000.000000000000000000, 1.000000000000000000, 50000.000000000000000000, 1.000000000000000000, 9000000.000000000000000000, 0.001200000000000000, 0.000000000000000000, 50000.000000000000000000, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, 0.000000000000000000, false, null, null, null, 'main', 50000.000000000000000000, 'a96c7467-c155-4de8-4da7-088ab56751fa');
		(2, 5, 'btcusdt', 'limit', 'sell', 'filled', '2021-03-17 12:44:36.481090', '2021-03-17 12:44:37.196980', 'none', 50000.000000000000000000, 1.000000000000000000, 1.000000000000000000, 1.000000000000000000, 10000000.000000000000000000, 125.000000000000000000, 0.000000000000000000, 1.000000000000000000, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, 0.000000000000000000, false, null, null, null, 'main', 50000.000000000000000000, 'acf782c3-5922-4345-7225-dbe1e623f15e');
		(3, 5, 'btcusdt', 'limit', 'sell', 'filled', '2021-03-17 12:45:46.879647', '2021-03-17 12:45:49.501344', 'none', 50000.000000000000000000, 1.000000000000000000, 1.000000000000000000, 1.000000000000000000, 9999999.998800000000000000, 60.000000000000000000, 0.000000000000000000, 1.000000000000000000, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, 0.000000000000000000, false, null, null, null, 'main', 50000.000000000000000000, '271f8de9-32d7-4667-4ec2-bd48c05945d0');
		(4, 5, 'btcusdt', 'limit', 'buy', 'filled', '2021-03-17 12:45:49.129942', '2021-03-17 12:45:49.502608', 'none', 50000.000000000000000000, 1.000000000000000000, 50000.000000000000000000, 1.000000000000000000, 8999875.000000000000000000, 0.002500000000000000, 0.000000000000000000, 50000.000000000000000000, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, 0.000000000000000000, false, null, null, null, 'main', 50000.000000000000000000, '519cc9af-c8cd-4e10-56e7-4467428f5fe6');
		(5, 5, 'btcusdt', 'limit', 'buy', 'filled', '2021-03-17 12:56:26.622327', '2021-03-17 12:56:30.098598', 'none', 50000.000000000000000000, 1.000000000000000000, 50000.000000000000000000, 1.000000000000000000, 8999815.000000000000000000, 0.001200000000000000, 0.000000000000000000, 50000.000000000000000000, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, 0.000000000000000000, false, null, null, null, 'main', 50000.000000000000000000, '4c5b2f04-19e2-46db-6bb5-a3d291fac469');
		(6, 5, 'btcusdt', 'limit', 'sell', 'filled', '2021-03-17 12:56:29.688058', '2021-03-17 12:56:30.097302', 'none', 50000.000000000000000000, 1.000000000000000000, 1.000000000000000000, 1.000000000000000000, 9999999.996300000000000000, 125.000000000000000000, 0.000000000000000000, 1.000000000000000000, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, 0.000000000000000000, false, null, null, null, 'main', 50000.000000000000000000, '7360e7a3-7b01-44b4-79d2-035767d3333a');
		(7, 5, 'btcusdt', 'limit', 'sell', 'filled', '2021-03-17 12:58:36.801136', '2021-03-17 12:58:49.183262', 'none', 50000.000000000000000000, 1.000000000000000000, 1.000000000000000000, 1.000000000000000000, 9999999.995100000000000000, 60.000000000000000000, 0.000000000000000000, 1.000000000000000000, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, 0.000000000000000000, false, null, null, null, 'main', 50000.000000000000000000, '357c3873-a7ba-469a-7bd4-9cb9d98dad99');
		(8, 5, 'btcusdt', 'limit', 'sell', 'untouched', '2021-03-17 12:58:42.958970', '2021-03-17 12:58:42.958970', 'none', 50000.000000000000000000, 1.000000000000000000, 1.000000000000000000, 0.000000000000000000, 9999998.995100000000000000, 0.000000000000000000, 0.000000000000000000, 0.000000000000000000, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, 0.000000000000000000, false, null, null, null, 'main', 0.000000000000000000, 'fa5dbb73-8628-427c-57b3-c291db878649');
		(9, 5, 'btcusdt', 'limit', 'buy', 'filled', '2021-03-17 12:58:48.289878', '2021-03-17 12:58:49.183737', 'none', 50000.000000000000000000, 1.000000000000000000, 50000.000000000000000000, 1.000000000000000000, 8999690.000000000000000000, 0.002500000000000000, 0.000000000000000000, 50000.000000000000000000, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, null, null, null, 0.000000000000000000, 0.000000000000000000, false, null, null, null, 'main', 50000.000000000000000000, 'd025269a-4e01-4254-7e54-76fe1011903f');
		*/

		Convey("Issue #003 [22.03.2021] - Duplicate statuses generated for limit orders", func() {
			orderBook := NewOrderBook("btcusdt", 2, 8)
			order7 := model.Order{ID: 7, Price: 5000000, Amount: 100000000, StopPrice: 0, Funds: 100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Limit, Side: model.MarketSide_Sell}
			order8 := model.Order{ID: 8, Price: 5000000, Amount: 100000000, StopPrice: 0, Funds: 100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Limit, Side: model.MarketSide_Sell}
			order9 := model.Order{ID: 9, Price: 5000000, Amount: 100000000, StopPrice: 0, Funds: 5000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Limit, Side: model.MarketSide_Buy}
			events = events[0:0]
			orderBook.Process(order7, &events)
			orderBook.Process(order8, &events)
			orderBook.Process(order9, &events)
			So(len(events), ShouldEqual, 6)
			So(events[5].GetOrderStatus().ID, ShouldEqual, 7)
			So(events[5].GetOrderStatus().Status, ShouldEqual, model.OrderStatus_Filled)
			So(events[5].GetOrderStatus().FilledAmount, ShouldEqual, 100000000)
			So(events[5].GetOrderStatus().UsedFunds, ShouldEqual, 5000000)
		})

	})
}
