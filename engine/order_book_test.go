package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOrderBookProcessing(t *testing.T) {
	book := NewOrderBook()
	Convey("Given an empty order book", t, func() {
		//	BUY          SELL
		//	1.0 120      -
		Convey("Add a first buy order", func() {
			trades := book.Process(NewOrder("TEST_1", uint64(100000000), uint64(12000000000), 1, 1, 1))
			So(len(trades), ShouldEqual, 0)
		})
		// BUY          SELL
		// 1.0 120      1.1 120
		Convey("Add a first sell order", func() {
			trades := book.Process(NewOrder("TEST_2", uint64(110000000), uint64(12000000000), 2, 1, 1))
			So(len(trades), ShouldEqual, 0)
		})
		// BUY          SELL
		// -            1.1 120
		Convey("Add a matching sell order", func() {
			trades := book.Process(NewOrder("TEST_3", uint64(90000000), uint64(12000000000), 2, 1, 1))
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
			trades := book.Process(NewOrder("TEST_4", uint64(110000000), uint64(2000000000), 1, 1, 1))
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
			trades := book.Process(NewOrder("TEST_5", uint64(111000000), uint64(12000000000), 1, 1, 1))
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
			trades := book.Process(NewOrder("TEST_6", uint64(120000000), uint64(12000000000), 1, 1, 1))
			So(len(trades), ShouldEqual, 0)
			trades = book.Process(NewOrder("TEST_6_1", uint64(120000000), uint64(1000000000), 1, 1, 1))
			So(len(trades), ShouldEqual, 0)
		})
		// ORDER: SELL 1.30 10
		// BUY          SELL
		// 1.20 130     1.30 10
		// 1.11 20
		Convey("Add two sell orders at the same price without matching", func() {
			trades := book.Process(NewOrder("TEST_7", uint64(130000000), uint64(1000000000), 2, 1, 1))
			So(len(trades), ShouldEqual, 0)
			trades = book.Process(NewOrder("TEST_7_1", uint64(130000000), uint64(1000000000), 2, 1, 1))
			So(len(trades), ShouldEqual, 0)
		})
		// ORDER: SELL 1.40 20
		// BUY          SELL
		// 1.20 130     -
		// 1.11 20
		Convey("Add a buy order that clears the sell side of the order book", func() {
			trades := book.Process(NewOrder("TEST_8", uint64(140000000), uint64(2000000000), 1, 1, 1))
			So(len(trades), ShouldEqual, 2)
		})
		// ORDER: SELL 1.00 140
		// BUY          SELL
		// -            -
		Convey("Add a sell order that clears the buy side of the order book", func() {
			trades := book.Process(NewOrder("TEST_9", uint64(100000000), uint64(15000000000), 2, 1, 1))
			So(len(trades), ShouldEqual, 3)

			state := book.GetMarket()
			So(state[0].Len(), ShouldEqual, 0)
			So(state[1].Len(), ShouldEqual, 0)
		})
		Convey("Lowest ask price should be updated", func() {
			book.Process(NewOrder("TEST_10", uint64(110000000), uint64(1000000000), 2, 1, 1))
			book.Process(NewOrder("TEST_11", uint64(130000000), uint64(1000000000), 2, 1, 1))
			trades := book.Process(NewOrder("TEST_12", uint64(140000000), uint64(1300000000), 1, 1, 1))
			So(len(trades), ShouldEqual, 2)
			So(book.GetLowestAsk(), ShouldEqual, uint64(130000000))
		})

		Convey("Highest ask should be 0 when there are no more buy orders", func() {
			book.Process(NewOrder("TEST_13", uint64(140000000), uint64(1000000000), 1, 1, 1))
			trades := book.Process(NewOrder("TEST_14", uint64(110000000), uint64(800000000), 2, 1, 1))
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
			book.Process(NewOrder("TB_1", uint64(111000000), uint64(100000000), 1, 1, 1))
			book.Process(NewOrder("TB_15", uint64(121000000), uint64(700000000), 1, 1, 1))

			//sell
			book.Process(NewOrder("TS_7", uint64(116000000), uint64(700000000), 2, 1, 1))
			book.Process(NewOrder("TS_8", uint64(113000000), uint64(700000000), 2, 1, 1))
			book.Process(NewOrder("TS_9", uint64(113000000), uint64(700000000), 2, 1, 1))

			backup := book.Backup()

			So(backup.HighestBid, ShouldBeLessThan, backup.LowestAsk)
		})

		Convey("Check if the lowest ask is moved after a completed limit order", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))
			orderBook.Process(NewOrder("mk_2", uint64(120000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))
			trades := orderBook.Process(NewOrder("mk_3", uint64(110000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("Check if the highest bid is moved after a completed limit order", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(120000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			orderBook.Process(NewOrder("mk_2", uint64(110000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			trades := orderBook.Process(NewOrder("mk_3", uint64(110000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("Check if trades are returned if the pricepoint is not complete", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(120000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			orderBook.Process(NewOrder("mk_2", uint64(110000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			trades := orderBook.Process(NewOrder("mk_3", uint64(110000000), uint64(11000000000), SELL, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))

			trades := orderBook.Process(Order{ID: "mk_2", Amount: 10000000000, Funds: 1100000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 1000000000)

			trades = orderBook.Process(Order{ID: "mk_3", Amount: 11000000000, Funds: 12100000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 2", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))

			trades := orderBook.Process(Order{ID: "mk_2", Amount: 10000000000, Funds: 1100000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 1000000000)

			trades = orderBook.Process(Order{ID: "mk_3", Amount: 15000000000, Funds: 23456000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 3", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))
			orderBook.Process(NewOrder("mk_2", uint64(120000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))

			trades := orderBook.Process(Order{ID: "mk_3", Amount: 10000000000, Funds: 1100000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 1000000000)

			trades = orderBook.Process(Order{ID: "mk_4", Amount: 15000000000, Funds: 23456000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a market buy order in an existing market 4", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))
			orderBook.Process(NewOrder("mk_2", uint64(120000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))

			trades := orderBook.Process(Order{ID: "mk_3", Amount: 12000000000, Funds: 51100000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			trades := orderBook.Process(Order{ID: "mk_2", Amount: 10000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 1)
			trades = orderBook.Process(Order{ID: "mk_3", Amount: 10000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 1)
		})

		Convey("I should be able to add a market sell order in an existing market 1", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_4", uint64(110000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			trades := orderBook.Process(Order{ID: "mk_5", Amount: 13000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 2", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))

			trades := orderBook.Process(Order{ID: "mk_2", Amount: 10000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 10000000000)

			trades = orderBook.Process(Order{ID: "mk_3", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 2000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 3", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))

			trades := orderBook.Process(Order{ID: "mk_2", Amount: 12000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 110000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 4", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(120000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))

			trades := orderBook.Process(Order{ID: "mk_2", Amount: 12000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 12000000000)
		})

		Convey("I should be able to add a market sell order in an existing market 5", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(120000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			orderBook.Process(NewOrder("mk_1", uint64(110000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))

			trades := orderBook.Process(Order{ID: "mk_2", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
			trades = orderBook.Process(Order{ID: "mk_3", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 2)
		})

		Convey("I should be able to add a market buy order in an empty market", func() {
			orderBook := NewOrderBook()
			trades := orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, Funds: 1000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			So(len(trades), ShouldEqual, 0)
		})

		Convey("I should be able to add a market sell order in an empty market", func() {
			orderBook := NewOrderBook()
			trades := orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 0)
		})

		Convey("I should be able to add a market buy order after another pending market order", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(120000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))
			orderBook.Process(Order{ID: "mk_2", Amount: 14000000000, Funds: 20000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			trades := orderBook.Process(Order{ID: "mk_3", Amount: 11000000000, Funds: 1000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			So(len(trades), ShouldEqual, 0)
		})

		Convey("I should be able to add a market sell order after another pending market order", func() {
			orderBook := NewOrderBook()
			orderBook.Process(NewOrder("mk_1", uint64(120000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))
			orderBook.Process(Order{ID: "mk_2", Amount: 13000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			trades := orderBook.Process(Order{ID: "mk_3", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			So(len(trades), ShouldEqual, 0)
		})

		Convey("I should be able to add a limit sell order after another pending market buy order", func() {
			orderBook := NewOrderBook()
			orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, Funds: 13200000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			trades := orderBook.Process(NewOrder("mk_2", uint64(120000000), uint64(12000000000), SELL, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit buy order after another pending market sell order", func() {
			orderBook := NewOrderBook()
			orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			trades := orderBook.Process(NewOrder("mk_2", uint64(120000000), uint64(12000000000), BUY, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders", func() {
			orderBook := NewOrderBook()
			orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, Funds: 13200000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			orderBook.Process(Order{ID: "mk_3", Amount: 11000000000, Funds: 13200000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			trades := orderBook.Process(NewOrder("mk_2", uint64(120000000), uint64(24000000000), SELL, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit buy order after multiple pending market sell orders", func() {
			orderBook := NewOrderBook()
			orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			orderBook.Process(Order{ID: "mk_2", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			trades := orderBook.Process(NewOrder("mk_3", uint64(120000000), uint64(24000000000), BUY, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders 2", func() {
			orderBook := NewOrderBook()
			orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, Funds: 13200000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			orderBook.Process(Order{ID: "mk_3", Amount: 20000000000, Funds: 50000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			trades := orderBook.Process(NewOrder("mk_2", uint64(120000000), uint64(24000000000), SELL, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)

			So(trades[1].Price, ShouldEqual, 120000000)
			So(trades[1].Amount, ShouldEqual, 13000000000)
		})

		Convey("I should be able to add a limit buy order after multiple pending market sell orders 2", func() {
			orderBook := NewOrderBook()
			orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			orderBook.Process(Order{ID: "mk_2", Amount: 20000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			trades := orderBook.Process(NewOrder("mk_3", uint64(120000000), uint64(24000000000), BUY, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
			So(trades[1].Price, ShouldEqual, 120000000)
			So(trades[1].Amount, ShouldEqual, 13000000000)
		})

		Convey("I should be able to add a limit sell order after multiple pending market buy orders 3", func() {
			orderBook := NewOrderBook()
			orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, Funds: 13200000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			orderBook.Process(Order{ID: "mk_2", Amount: 20000000000, Funds: 50000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY})
			trades := orderBook.Process(NewOrder("mk_3", uint64(120000000), uint64(24000000000), SELL, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
			So(trades[1].Price, ShouldEqual, 120000000)
			So(trades[1].Amount, ShouldEqual, 13000000000)
			trades = orderBook.Process(NewOrder("mk_4", uint64(120000000), uint64(24000000000), SELL, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Amount, ShouldEqual, 7000000000)
		})

		Convey("I should be able to add a limit buy order after multiple pending market sell orders 3", func() {
			orderBook := NewOrderBook()
			orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			orderBook.Process(Order{ID: "mk_2", Amount: 20000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL})
			trades := orderBook.Process(NewOrder("mk_3", uint64(120000000), uint64(24000000000), BUY, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 2)
			So(trades[0].Price, ShouldEqual, 120000000)
			So(trades[0].Amount, ShouldEqual, 11000000000)
			So(trades[1].Amount, ShouldEqual, 13000000000)
			trades = orderBook.Process(NewOrder("mk_4", uint64(120000000), uint64(24000000000), BUY, LimitOrder, EventTypeNewOrder))
			So(len(trades), ShouldEqual, 1)
			So(trades[0].Amount, ShouldEqual, 7000000000)
		})

		Convey("Should not process invalid order types", func() {
			orderBook := NewOrderBook()
			trades := orderBook.Process(Order{ID: "mk_1", Amount: 11000000000, EventType: 0, Type: MarketOrder, Side: SELL})
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
			orderBook := NewOrderBook()
			order := NewOrder("TEST_15", uint64(110000000), uint64(800000000), 2, 1, 1)
			orderBook.Process(order)
			So(orderBook.Cancel(order), ShouldBeTrue)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
		})

		Convey("I should be able to cancel a buy order", func() {
			orderBook := NewOrderBook()
			order := NewOrder("TEST_15", uint64(110000000), uint64(800000000), 1, 1, 1)
			orderBook.Process(order)
			So(orderBook.Cancel(order), ShouldBeTrue)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
		})

		Convey("I should be able to cancel a buy market order", func() {
			orderBook := NewOrderBook()
			order := Order{ID: "mk_1", Amount: 11000000000, Funds: 13200000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY}
			orderBook.Process(order)
			So(orderBook.Cancel(order), ShouldBeTrue)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
		})

		Convey("I should be able to cancel a sell market order", func() {
			orderBook := NewOrderBook()
			order := Order{ID: "mk_1", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL}
			orderBook.Process(order)
			So(orderBook.Cancel(order), ShouldBeTrue)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)
		})

		Convey("I should be able to cancel an invalid order", func() {
			orderBook := NewOrderBook()
			order := NewOrder("TEST_15", uint64(110000000), uint64(800000000), 1, 1, 1)
			So(orderBook.Cancel(order), ShouldBeFalse)

			state := orderBook.GetMarket()
			So(state[0].Len(), ShouldBeZeroValue)
			So(state[1].Len(), ShouldBeZeroValue)

			order = Order{ID: "mk_1", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: SELL}
			So(orderBook.Cancel(order), ShouldBeFalse)

			order = Order{ID: "mk_1", Amount: 11000000000, EventType: EventTypeNewOrder, Type: MarketOrder, Side: BUY}
			So(orderBook.Cancel(order), ShouldBeFalse)

			order = Order{ID: "mk_1", Amount: 11000000000, EventType: EventTypeNewOrder, Type: 0, Side: BUY}
			So(orderBook.Cancel(order), ShouldBeFalse)
		})

	})
}
