package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOrderBookProcessing(t *testing.T) {
	book := NewOrderBook("btcusd", 8, 8)
	trades := make([]Trade, 0, 5)
	Convey("Given an empty order book", t, func() {
		trades = trades[0:0]
		//	BUY          SELL
		//	1.0 120      -
		Convey("Add a first buy order", func() {
			book.Process(NewOrder(1, uint64(100000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 0)
		})
		// BUY          SELL
		// 1.0 120      1.1 120
		Convey("Add a first sell order", func() {
			book.Process(NewOrder(2, uint64(110000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 0)
		})
		// BUY          SELL
		// -            1.1 120
		Convey("Add a matching sell order", func() {
			book.Process(NewOrder(3, uint64(90000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			trade := trades[0]
			So(trade.Amount, ShouldEqual, 12000000000)
			So(trade.Price, ShouldEqual, 100000000)

			state := book.GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 1)
		})
		// ORDER: BUY 1.1 20
		// BUY          SELL
		// -            1.1 100
		Convey("Add a sell order with the same price", func() {
			book.Process(NewOrder(4, uint64(110000000), uint64(2000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			trade := trades[0]
			So(trade.Amount, ShouldEqual, 2000000000)
			So(trade.Price, ShouldEqual, 110000000)

			So(book.GetLowestAsk(), ShouldEqual, 110000000)
		})
		// ORDER: BUY 1.11 120
		// BUY          SELL
		// 1.11 20      -
		Convey("Add a buy order with a larger amount than the available sell", func() {
			book.Process(NewOrder(5, uint64(111000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			trade := trades[0]
			So(trade.Amount, ShouldEqual, 10000000000)
			So(trade.Price, ShouldEqual, 110000000)
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
			book.Process(NewOrder(6, uint64(120000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 0)
			trades = trades[0:0]
			book.Process(NewOrder(7, uint64(120000000), uint64(1000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 0)
		})
		// ORDER: SELL 1.30 10
		// BUY          SELL
		// 1.20 130     1.30 10
		// 1.11 20
		Convey("Add two sell orders at the same price without matching", func() {
			book.Process(NewOrder(7, uint64(130000000), uint64(1000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 0)
			trades = trades[0:0]
			book.Process(NewOrder(8, uint64(130000000), uint64(1000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 0)
		})
		// ORDER: SELL 1.40 20
		// BUY          SELL
		// 1.20 130     -
		// 1.11 20
		Convey("Add a buy order that clears the sell side of the order book", func() {
			book.Process(NewOrder(9, uint64(140000000), uint64(2000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 2)
		})
		// ORDER: SELL 1.00 140
		// BUY          SELL
		// -            -
		Convey("Add a sell order that clears the buy side of the order book", func() {
			book.Process(NewOrder(10, uint64(100000000), uint64(15000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 3)

			state := book.GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 0)
		})
		Convey("Lowest ask price should be updated", func() {
			book.Process(NewOrder(11, uint64(110000000), uint64(1000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			book.Process(NewOrder(12, uint64(130000000), uint64(1000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			trades = trades[0:0]
			book.Process(NewOrder(13, uint64(140000000), uint64(1300000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 2)
			So(book.GetLowestAsk(), ShouldEqual, uint64(130000000))
		})

		Convey("Highest ask should be 0 when there are no more buy orders", func() {
			book.Process(NewOrder(14, uint64(140000000), uint64(1000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			trades = trades[0:0]
			book.Process(NewOrder(15, uint64(110000000), uint64(800000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			So(book.GetHighestBid(), ShouldEqual, 0)
		})

		Convey("Check if adding multiple orders works as expected", func() {
			/*
			   BUY        SELL
			   1.21 70    1.13 140 -> 1.11 10     1.13 70
			   1.11 10    1.16 70
			*/

			//buy
			book.Process(NewOrder(16, uint64(111000000), uint64(100000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			book.Process(NewOrder(17, uint64(121000000), uint64(700000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)

			//sell
			book.Process(NewOrder(18, uint64(116000000), uint64(700000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			book.Process(NewOrder(19, uint64(113000000), uint64(700000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			book.Process(NewOrder(20, uint64(113000000), uint64(700000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)

			backup := book.Backup()

			So(backup.HighestBid, ShouldBeLessThan, backup.LowestAsk)
		})

		Convey("Check if the lowest ask is moved after a completed limit order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(20, uint64(110000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			orderBook.Process(NewOrder(21, uint64(120000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(22, uint64(110000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("Check if the highest bid is moved after a completed limit order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(23, uint64(120000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			orderBook.Process(NewOrder(24, uint64(110000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(25, uint64(110000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("Check if trades are returned if the pricepoint is not complete", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(26, uint64(120000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			orderBook.Process(NewOrder(27, uint64(110000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(28, uint64(110000000), uint64(11000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(29, uint64(110000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 30, Amount: 10000000000, Funds: 1100000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 1000000000)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 31, Amount: 11000000000, Funds: 12100000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 2", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(32, uint64(110000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 33, Amount: 10000000000, Funds: 1100000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 1000000000)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 34, Amount: 15000000000, Funds: 23456000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 3", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(35, uint64(110000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			orderBook.Process(NewOrder(36, uint64(120000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 37, Amount: 10000000000, Funds: 1100000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 1000000000)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 38, Amount: 15000000000, Funds: 23456000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 4", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(39, uint64(110000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			orderBook.Process(NewOrder(40, uint64(120000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 41, Amount: 12000000000, Funds: 51100000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(42, uint64(110000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			trades = trades[0:0]
			orderBook.Process(Order{ID: 43, Amount: 10000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 1)
			trades = trades[0:0]
			orderBook.Process(Order{ID: 44, Amount: 10000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 1)
		})

		Convey("I should be able to add a market sell order in an existing market 1", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(45, uint64(110000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			trades = trades[0:0]
			orderBook.Process(Order{ID: 46, Amount: 13000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 2", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(47, uint64(110000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 48, Amount: 10000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 10000000000)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 49, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 2000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 3", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(50, uint64(110000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 51, Amount: 12000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 4", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(52, uint64(120000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			orderBook.Process(NewOrder(53, uint64(110000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 54, Amount: 12000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 5", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(55, uint64(120000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			orderBook.Process(NewOrder(56, uint64(110000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)

			trades = trades[0:0]
			orderBook.Process(Order{ID: 57, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
			trades = trades[0:0]
			orderBook.Process(Order{ID: 58, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 2)
		})

		Convey("I should be able to add a market buy order in an empty market", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 59, Amount: 11000000000, Funds: 1000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			So(len(trades), ShouldEqual, 0)
		})

		Convey("I should be able to add a market sell order in an empty market", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 60, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 0)
		})

		Convey("I should be able to add a market buy order after another pending market order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(61, uint64(120000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			orderBook.Process(Order{ID: 62, Amount: 14000000000, Funds: 20000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			trades = trades[0:0]
			orderBook.Process(Order{ID: 63, Amount: 11000000000, Funds: 1000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			So(len(trades), ShouldEqual, 0)
		})

		Convey("I should be able to add a market sell order after another pending market order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(NewOrder(64, uint64(120000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			orderBook.Process(Order{ID: 65, Amount: 13000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			trades = trades[0:0]
			orderBook.Process(Order{ID: 66, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 0)
		})

		Convey("I should be able to add a limit sell order after another pending market buy order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 67, Amount: 11000000000, Funds: 13200000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(68, uint64(120000000), uint64(12000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit buy order after another pending market sell order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 69, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(70, uint64(120000000), uint64(12000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 71, Amount: 11000000000, Funds: 13200000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			orderBook.Process(Order{ID: 72, Amount: 11000000000, Funds: 13200000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(73, uint64(120000000), uint64(24000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit buy order after multiple pending market sell orders", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 74, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			orderBook.Process(Order{ID: 75, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(76, uint64(120000000), uint64(24000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders 2", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 77, Amount: 11000000000, Funds: 13200000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			orderBook.Process(Order{ID: 78, Amount: 20000000000, Funds: 50000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(79, uint64(120000000), uint64(24000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)

			So(trades[1].Price, ShouldEqual, 120000000)
			So(trades[1].Amount, ShouldEqual, 13000000000)
		})

		Convey("I should be able to add a limit buy order after multiple pending market sell orders 2", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 80, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			orderBook.Process(Order{ID: 81, Amount: 20000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(82, uint64(120000000), uint64(24000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
			So(trades[1].Price, ShouldEqual, 120000000)
			So(trades[1].Amount, ShouldEqual, 13000000000)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders 3", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 83, Amount: 11000000000, Funds: 13200000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			orderBook.Process(Order{ID: 84, Amount: 20000000000, Funds: 50000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}, &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(85, uint64(120000000), uint64(24000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
			So(trades[1].Price, ShouldEqual, 120000000)
			So(trades[1].Amount, ShouldEqual, 13000000000)
			trades = trades[0:0]
			orderBook.Process(NewOrder(86, uint64(120000000), uint64(24000000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Amount, ShouldEqual, 7000000000)
		})

		Convey("I should be able to add a limit buy order after multiple pending market sell orders 3", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			orderBook.Process(Order{ID: 87, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			orderBook.Process(Order{ID: 88, Amount: 20000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			trades = trades[0:0]
			orderBook.Process(NewOrder(89, uint64(120000000), uint64(24000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
			So(trades[1].Amount, ShouldEqual, 13000000000)
			trades = trades[0:0]
			orderBook.Process(NewOrder(90, uint64(120000000), uint64(24000000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder), &trades)
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Amount, ShouldEqual, 7000000000)
		})

		Convey("Should not process invalid order types", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			trades = trades[0:0]
			orderBook.Process(Order{ID: 91, Amount: 11000000000, EventType: 0, Type: OrderType_Market, Side: MarketSide_Sell}, &trades)
			So(len(trades), ShouldEqual, 0)
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
			order := NewOrder(92, uint64(110000000), uint64(800000000), MarketSide_Sell, OrderType_Limit, CommandType_NewOrder)
			orderBook.Process(order, &trades)
			So(orderBook.Cancel(order), ShouldBeTrue)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
		})

		Convey("I should be able to cancel a buy order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			order := NewOrder(93, uint64(110000000), uint64(800000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder)
			orderBook.Process(order, &trades)
			So(orderBook.Cancel(order), ShouldBeTrue)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
		})

		Convey("I should be able to cancel a buy market order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			order := Order{ID: 94, Amount: 11000000000, Funds: 13200000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}
			orderBook.Process(order, &trades)
			So(orderBook.Cancel(order), ShouldBeTrue)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
		})

		Convey("I should be able to cancel a sell market order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			order := Order{ID: 95, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}
			orderBook.Process(order, &trades)
			So(orderBook.Cancel(order), ShouldBeTrue)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
		})

		Convey("I should be able to cancel an invalid order", func() {
			orderBook := NewOrderBook("btcusd", 8, 8)
			order := NewOrder(96, uint64(110000000), uint64(800000000), MarketSide_Buy, OrderType_Limit, CommandType_NewOrder)
			So(orderBook.Cancel(order), ShouldBeFalse)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)

			order = Order{ID: 97, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Sell}
			So(orderBook.Cancel(order), ShouldBeFalse)

			order = Order{ID: 98, Amount: 11000000000, EventType: CommandType_NewOrder, Type: OrderType_Market, Side: MarketSide_Buy}
			So(orderBook.Cancel(order), ShouldBeFalse)

			order = Order{ID: 99, Amount: 11000000000, EventType: CommandType_NewOrder, Type: 4, Side: MarketSide_Buy}
			So(orderBook.Cancel(order), ShouldBeFalse)
		})

	})
}
