package trading_engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBookEntryCreation(t *testing.T) {
	Convey("Given an order", t, func() {
		order := NewOrder("TEST_1", 100000000, 12000000000, 1, 1)
		Convey("I should be able to create a new order book", func() {
			bookEntry := NewBookEntry(order)
			So(bookEntry.Amount, ShouldEqual, 12000000000)
			So(bookEntry.Price, ShouldEqual, 100000000)
			So(bookEntry.Order, ShouldResemble, order)
		})
	})
}
