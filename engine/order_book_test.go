package engine

import (
	"testing"

	"gitlab.com/around25/products/matching-engine/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOrderBookProcessing(t *testing.T) {
	book := NewOrderBook("btcusd", 8, 8)
	events := make([]model.Event, 0, 5)
	Convey("Given an empty order book", t, func() {
		events = events[0:0]
		//	BUY          SELL
		//	1.0 120      -
		Convey("Add a first buy order", func() {
			book.Process(model.NewOrder(1, uint64(100000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 0)
		})
		// BUY          SELL
		// 1.0 120      1.1 120
		Convey("Add a first sell order", func() {
			book.Process(model.NewOrder(2, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 0)
		})
		// BUY          SELL
		// -            1.1 120
		Convey("Add a matching sell order", func() {
			book.Process(model.NewOrder(3, uint64(90000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Amount, ShouldEqual, 12000000000)
			So(events[0].GetTrade().Price, ShouldEqual, 100000000)

			state := book.GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 1)
		})
		// ORDER: BUY 1.1 20
		// BUY          SELL
		// -            1.1 100
		Convey("Add a sell order with the same price", func() {
			book.Process(model.NewOrder(4, uint64(110000000), uint64(2000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Amount, ShouldEqual, 2000000000)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)

			So(book.GetLowestAsk(), ShouldEqual, 110000000)
		})
		// ORDER: BUY 1.11 120
		// BUY          SELL
		// 1.11 20      -
		Convey("Add a buy order with a larger amount than the available sell", func() {
			book.Process(model.NewOrder(5, uint64(111000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Amount, ShouldEqual, 10000000000)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
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
			So(len(events), ShouldEqual, 0)
			events = events[0:0]
			book.Process(model.NewOrder(7, uint64(120000000), uint64(1000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 0)
		})
		// ORDER: SELL 1.30 10
		// BUY          SELL
		// 1.20 130     1.30 10
		// 1.11 20
		Convey("Add two sell orders at the same price without matching", func() {
			book.Process(model.NewOrder(7, uint64(130000000), uint64(1000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 0)
			events = events[0:0]
			book.Process(model.NewOrder(8, uint64(130000000), uint64(1000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 0)
		})
		// ORDER: SELL 1.40 20
		// BUY          SELL
		// 1.20 130     -
		// 1.11 20
		Convey("Add a buy order that clears the sell side of the order book", func() {
			book.Process(model.NewOrder(9, uint64(140000000), uint64(2000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 2)
		})
		// ORDER: SELL 1.00 140
		// BUY          SELL
		// -            -
		Convey("Add a sell order that clears the buy side of the order book", func() {
			book.Process(model.NewOrder(10, uint64(100000000), uint64(15000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 3)

			state := book.GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 0)
		})
		Convey("Lowest ask price should be updated", func() {
			book.Process(model.NewOrder(11, uint64(110000000), uint64(1000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			book.Process(model.NewOrder(12, uint64(130000000), uint64(1000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			book.Process(model.NewOrder(13, uint64(140000000), uint64(1300000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 2)
			So(book.GetLowestAsk(), ShouldEqual, uint64(130000000))
		})

		Convey("Highest ask should be 0 when there are no more buy orders", func() {
			book.Process(model.NewOrder(14, uint64(140000000), uint64(1000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			book.Process(model.NewOrder(15, uint64(110000000), uint64(800000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
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
			orderBook.Process(model.NewOrder(20, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(21, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(22, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 12000000000)
		})

		Convey("Check if the highest bid is moved after a completed limit order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(23, uint64(120000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(24, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(25, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 12000000000)
		})

		Convey("Check if events are returned if the pricepoint is not complete", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(26, uint64(120000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(27, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(28, uint64(110000000), uint64(11000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(29, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 30, Amount: 10000000000, Funds: 1100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 1000000000)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 31, Amount: 11000000000, Funds: 12100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 2", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(32, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 33, Amount: 10000000000, Funds: 1100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 1000000000)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 34, Amount: 15000000000, Funds: 23456000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 3", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(35, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(36, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 37, Amount: 10000000000, Funds: 1100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 1000000000)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 38, Amount: 15000000000, Funds: 23456000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 2)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 4", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(39, uint64(110000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(40, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 41, Amount: 12000000000, Funds: 51100000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(42, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			orderBook.Process(model.Order{ID: 43, Amount: 10000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 1)
			events = events[0:0]
			orderBook.Process(model.Order{ID: 44, Amount: 10000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 1)
		})

		Convey("I should be able to add a market sell order in an existing market 1", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(45, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			events = events[0:0]
			orderBook.Process(model.Order{ID: 46, Amount: 13000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 2", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(47, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 48, Amount: 10000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 10000000000)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 49, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 2000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 3", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(50, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 51, Amount: 12000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 110000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 4", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(52, uint64(120000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(53, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 54, Amount: 12000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 5", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(55, uint64(120000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.NewOrder(56, uint64(110000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)

			events = events[0:0]
			orderBook.Process(model.Order{ID: 57, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
			events = events[0:0]
			orderBook.Process(model.Order{ID: 58, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 2)
		})

		Convey("I should be able to add a market buy order in an empty market", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 59, Amount: 11000000000, Funds: 1000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 0)
		})

		Convey("I should be able to add a market sell order in an empty market", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 60, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 0)
		})

		Convey("I should be able to add a market buy order after another pending market order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(61, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.Order{ID: 62, Amount: 14000000000, Funds: 20000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			events = events[0:0]
			orderBook.Process(model.Order{ID: 63, Amount: 11000000000, Funds: 1000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			So(len(events), ShouldEqual, 0)
		})

		Convey("I should be able to add a market sell order after another pending market order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.NewOrder(64, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			orderBook.Process(model.Order{ID: 65, Amount: 13000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			events = events[0:0]
			orderBook.Process(model.Order{ID: 66, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 0)
		})

		Convey("I should be able to add a limit sell order after another pending market buy order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 67, Amount: 11000000000, Funds: 13200000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(68, uint64(120000000), uint64(12000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit buy order after another pending market sell order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 69, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(70, uint64(120000000), uint64(12000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 71, Amount: 11000000000, Funds: 13200000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			orderBook.Process(model.Order{ID: 72, Amount: 11000000000, Funds: 13200000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(73, uint64(120000000), uint64(24000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 2)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit buy order after multiple pending market sell orders", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 74, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			orderBook.Process(model.Order{ID: 75, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(76, uint64(120000000), uint64(24000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 2)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders 2", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 77, Amount: 11000000000, Funds: 13200000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			orderBook.Process(model.Order{ID: 78, Amount: 20000000000, Funds: 50000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(79, uint64(120000000), uint64(24000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 2)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)

			So(events[1].GetTrade().Price, ShouldEqual, 120000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 13000000000)
		})

		Convey("I should be able to add a limit buy order after multiple pending market sell orders 2", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 80, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			orderBook.Process(model.Order{ID: 81, Amount: 20000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(82, uint64(120000000), uint64(24000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 2)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
			So(events[1].GetTrade().Price, ShouldEqual, 120000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 13000000000)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders 3", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 83, Amount: 11000000000, Funds: 13200000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			orderBook.Process(model.Order{ID: 84, Amount: 20000000000, Funds: 50000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(85, uint64(120000000), uint64(24000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 2)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
			So(events[1].GetTrade().Price, ShouldEqual, 120000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 13000000000)
			events = events[0:0]
			orderBook.Process(model.NewOrder(86, uint64(120000000), uint64(24000000000), model.MarketSide_Sell, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Amount, ShouldEqual, 7000000000)
		})

		Convey("I should be able to add a limit buy order after multiple pending market sell orders 3", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(model.Order{ID: 87, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			orderBook.Process(model.Order{ID: 88, Amount: 20000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			events = events[0:0]
			orderBook.Process(model.NewOrder(89, uint64(120000000), uint64(24000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 2)
			So(events[0].GetTrade().Price, ShouldEqual, 120000000)
			So(events[0].GetTrade().Amount, ShouldEqual, 11000000000)
			So(events[1].GetTrade().Amount, ShouldEqual, 13000000000)
			events = events[0:0]
			orderBook.Process(model.NewOrder(90, uint64(120000000), uint64(24000000000), model.MarketSide_Buy, model.OrderType_Limit, model.CommandType_NewOrder), &events)
			So(len(events), ShouldEqual, 1)
			So(events[0].GetTrade().Amount, ShouldEqual, 7000000000)
		})

		Convey("Should not process invalid order types", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			events = events[0:0]
			orderBook.Process(model.Order{ID: 91, Amount: 11000000000, EventType: 0, Type: model.OrderType_Market, Side: model.MarketSide_Sell}, &events)
			So(len(events), ShouldEqual, 0)
		})

		Convey("I should receive null when poping an empty pending market orders list", func() {
			orderBook := orderBook{}
			buyOrder := orderBook.popMarketBuyOrder()
			sellOrder := orderBook.popMarketSellOrder()
			So(buyOrder, ShouldBeNil)
			So(sellOrder, ShouldBeNil)
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

		Convey("I should be able to cancel a buy market order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			order := model.Order{ID: 94, Amount: 11000000000, Funds: 13200000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Buy}
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

		Convey("I should be able to cancel a sell market order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			order := model.Order{ID: 95, Amount: 11000000000, EventType: model.CommandType_NewOrder, Type: model.OrderType_Market, Side: model.MarketSide_Sell}
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
			So(len(events), ShouldEqual, 0)

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
