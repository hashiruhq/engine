package engine_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/around25/products/matching-engine/engine"
)

func TestBookEntryCreation(t *testing.T) {
	Convey("Given an order", t, func() {
		order := engine.NewOrder("TEST_1", 100000000, 12000000000, 1, 1, 1)
		Convey("I should be able to create a new book entry", func() {
			bookEntry := engine.BookEntry{Price: order.Price, Amount: order.Amount, Order: &order}
			So(bookEntry.Amount, ShouldEqual, 12000000000)
			So(bookEntry.Price, ShouldEqual, 100000000)
			So(bookEntry.Order, ShouldEqual, &order)
		})
	})
}

func TestBookEntryLoadFromToBinary(t *testing.T) {
	Convey("Should be able to load a book entry from binary", t, func() {
		order := engine.Order{ID: "TST_1", Price: 131221300010201, Amount: 848382829993942, Market: "btc-usd", Funds: 10100010133232313}
		entry := engine.BookEntry{Price: 131221300010201, Amount: 848382829993942, Order: &order}
		encoded, _ := entry.ToBinary()
		entry.FromBinary(encoded)
		So((*engine.BookEntry)(nil).GetPrice(), ShouldBeZeroValue)
		So((*engine.BookEntry)(nil).GetAmount(), ShouldBeZeroValue)
		So((*engine.BookEntry)(nil).GetOrder(), ShouldBeZeroValue)
		So(entry.GetPrice(), ShouldEqual, 131221300010201)
		So(entry.GetAmount(), ShouldEqual, 848382829993942)
		So(entry.GetOrder().GetPrice(), ShouldEqual, 131221300010201)
		So(entry.GetOrder().GetAmount(), ShouldEqual, 848382829993942)
		So(entry.GetOrder().GetFunds(), ShouldEqual, 10100010133232313)
		So(entry.GetOrder().GetID(), ShouldEqual, "TST_1")
		So(entry.GetOrder().GetSide(), ShouldEqual, engine.MarketSide_Buy)
		So(entry.GetOrder().GetType(), ShouldEqual, engine.CommandType_NewOrder)
		So(entry.GetOrder().GetMarket(), ShouldEqual, "btc-usd")

		So(entry.String(), ShouldEqual, `Price:131221300010201 Amount:848382829993942 Order:<Amount:848382829993942 Price:131221300010201 Funds:10100010133232313 ID:"TST_1" Market:"btc-usd" > `)
		entry.ProtoMessage()
		fd, i := entry.Descriptor()
		So(fd, ShouldNotBeNil)
		So(i, ShouldNotBeNil)
	})
}
