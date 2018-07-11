package trading_engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPricePointCreation(t *testing.T) {
	Convey("I should be able to create Price Points list", t, func() {
		pricePoints := NewPricePoints()
		Convey("And add a new price point to it with a new book", func() {
			order := NewOrder("TEST_1", uint64(100000000), uint64(12000000000), 1, 1, 1)
			bookEntry := NewBookEntry(order)
			pricePoints.Set(bookEntry.Price, &PricePoint{
				BuyBookEntries: []BookEntry{bookEntry},
			})
			So(pricePoints.Len(), ShouldEqual, 1)
			Convey("And the order added should be the same as the one retrieved", func() {
				value, ok := pricePoints.Get(uint64(100000000))
				So(ok, ShouldBeTrue)
				entry := value.BuyBookEntries[0]
				So(entry.Price, ShouldEqual, 100000000)
			})
		})
	})
}
