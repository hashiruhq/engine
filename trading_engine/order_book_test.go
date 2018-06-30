package trading_engine

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
			trades := book.Process(NewOrder("TEST_1", uint64(100000000), uint64(12000000000), 1, 1))
			So(len(trades), ShouldEqual, 0)
		})
		// BUY          SELL
		// 1.0 120      1.1 120
		Convey("Add a first sell order", func() {
			trades := book.Process(NewOrder("TEST_2", uint64(110000000), uint64(12000000000), 2, 1))
			So(len(trades), ShouldEqual, 0)
		})
		// BUY          SELL
		// -            1.1 120
		Convey("Add a matching sell order", func() {
			trades := book.Process(NewOrder("TEST_3", uint64(90000000), uint64(12000000000), 2, 1))
			So(len(trades), ShouldEqual, 1)
			trade := trades[0]
			So(trade.Amount, ShouldEqual, 12000000000)
			So(trade.Price, ShouldEqual, 100000000)
			So(book.GetMarket().Len(), ShouldEqual, 1)
		})
		// ORDER: BUY 1.1 20
		// BUY          SELL
		// -            1.1 100
		Convey("Add a sell order with the same price", func() {
			trades := book.Process(NewOrder("TEST_4", uint64(110000000), uint64(2000000000), 1, 1))
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
			trades := book.Process(NewOrder("TEST_5", uint64(111000000), uint64(12000000000), 1, 1))
			So(len(trades), ShouldEqual, 1)
			trade := trades[0]
			So(trade.Amount, ShouldEqual, 10000000000)
			So(trade.Price, ShouldEqual, 110000000)
			So(book.GetHighestBid(), ShouldEqual, 111000000)
			So(book.GetMarket().Len(), ShouldEqual, 1)
		})
		// ORDER: BUY 1.20 120
		// BUY          SELL
		// 1.20 120      -
		// 1.11 20
		Convey("Add two another buy orders with a higher price", func() {
			trades := book.Process(NewOrder("TEST_6", uint64(120000000), uint64(12000000000), 1, 1))
			So(len(trades), ShouldEqual, 0)
			trades = book.Process(NewOrder("TEST_6_1", uint64(120000000), uint64(1000000000), 1, 1))
			So(len(trades), ShouldEqual, 0)
		})
		// ORDER: SELL 1.30 10
		// BUY          SELL
		// 1.20 130     1.30 10
		// 1.11 20
		Convey("Add two sell orders at the same price without matching", func() {
			trades := book.Process(NewOrder("TEST_7", uint64(130000000), uint64(1000000000), 2, 1))
			So(len(trades), ShouldEqual, 0)
			trades = book.Process(NewOrder("TEST_7_1", uint64(130000000), uint64(1000000000), 2, 1))
			So(len(trades), ShouldEqual, 0)
		})
		// ORDER: SELL 1.40 20
		// BUY          SELL
		// 1.20 130     -
		// 1.11 20
		Convey("Add a buy order that clears the sell side of the order book", func() {
			trades := book.Process(NewOrder("TEST_8", uint64(140000000), uint64(2000000000), 1, 1))
			So(len(trades), ShouldEqual, 2)
		})
		// ORDER: SELL 1.00 140
		// BUY          SELL
		// -            -
		Convey("Add a sell order that clears the buy side of the order book", func() {
			trades := book.Process(NewOrder("TEST_9", uint64(100000000), uint64(15000000000), 2, 1))
			So(len(trades), ShouldEqual, 3)
			So(book.GetMarket().Len(), ShouldEqual, 0)
		})
		Convey("Lowest ask price should be updated", func() {
			book.Process(NewOrder("TEST_10", uint64(110000000), uint64(1000000000), 2, 1))
			book.Process(NewOrder("TEST_11", uint64(130000000), uint64(1000000000), 2, 1))
			trades := book.Process(NewOrder("TEST_12", uint64(140000000), uint64(1300000000), 1, 1))
			So(len(trades), ShouldEqual, 2)
			So(book.GetLowestAsk(), ShouldEqual, uint64(130000000))
		})

		Convey("Highest ask should be 0 when there are no more buy orders", func() {
			book.Process(NewOrder("TEST_13", uint64(140000000), uint64(1000000000), 1, 1))
			trades := book.Process(NewOrder("TEST_14", uint64(110000000), uint64(800000000), 2, 1))
			So(len(trades), ShouldEqual, 1)
			So(book.GetHighestBid(), ShouldEqual, 0)
		})

		Convey("I should be able to cancel a sell order", func() {
			err := book.Cancel("TEST_14")
			So(err, ShouldBeNil)
			SkipSo(book.GetMarket().Len(), ShouldBeZeroValue)
		})
	})
}
