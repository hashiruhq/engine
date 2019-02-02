package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPricePointCreation(t *testing.T) {
	Convey("I should be able to create Price Points list", t, func() {
		pricePoints := NewPricePoints()
		Convey("And add a new price point to it with a new book", func() {
			order := NewOrder("TEST_1", uint64(100000000), uint64(12000000000), 1, 1, 1)
			bookEntry := BookEntry{Price: order.Price, Amount: order.Amount, Order: &order}
			pricePoints.Set(bookEntry.Price, &PricePoint{
				Entries: []BookEntry{bookEntry},
			})
			So(pricePoints.Len(), ShouldEqual, 1)
			Convey("And the order added should be the same as the one retrieved", func() {
				value, ok := pricePoints.Get(uint64(100000000))
				So(ok, ShouldBeTrue)
				entry := value.Entries[0]
				So(entry.Price, ShouldEqual, 100000000)
			})
		})
	})

	Convey("I should be able to create Price Points list and get an iterator for it", t, func() {
		pricePoints := NewPricePoints()
		// test iterator
		iterator := pricePoints.Iterator()
		So(iterator.Value(), ShouldBeNil)
		iterator.Close()
		// test with list items
		firstPoint := &PricePoint{Entries: []BookEntry{BookEntry{Amount: 1}}}
		secondPoint := &PricePoint{Entries: []BookEntry{BookEntry{Amount: 2}}}
		thirdPoint := &PricePoint{Entries: []BookEntry{BookEntry{Amount: 3}}}
		pricePoints.Set(10, firstPoint)
		pricePoints.Set(13, secondPoint)
		pricePoints.Set(15, thirdPoint)
		it := pricePoints.Iterator()
		it.Seek(10)
		So(it.Value(), ShouldEqual, firstPoint)
		it.Next()
		So(it.Value(), ShouldEqual, secondPoint)
		it.Next()
		it.Seek(10)
		So(it.Value(), ShouldEqual, firstPoint)
		it.Close()

		// range iterator
		rit := pricePoints.Range(13, 16)
		rit.Next()
		So(rit.Value(), ShouldEqual, secondPoint)
		rit.Next()
		So(rit.Value(), ShouldEqual, thirdPoint)
		rit.Previous()
		So(rit.Value(), ShouldEqual, secondPoint)
		rit.Seek(13)
		So(rit.Value(), ShouldEqual, secondPoint)
		rit.Previous()
		rit.Seek(20)
		rit.Seek(9)
		rit.Seek(13)
		rit.Close()

		// test if the first element does not have a previous
		rit = pricePoints.Range(10, 16)
		rit.Next()
		rit.Previous()
		rit.Seek(15)
		rit.Next()
		rit.Close()

		rit = pricePoints.Range(10, 15)
		rit.Seek(13)
		rit.Next()
		rit.Close()

		// test seek
		it = pricePoints.SeekToFirst()
		So(it.Value(), ShouldEqual, firstPoint)
		it.Next()
		So(it.Value(), ShouldEqual, secondPoint)
		it.Close()
		it = pricePoints.SeekToLast()
		So(it.Value(), ShouldEqual, thirdPoint)
		it.Previous()
		So(it.Value(), ShouldEqual, secondPoint)
		it.Close()

		// change second point
		secondPoint = &PricePoint{}
		pricePoints.Set(13, secondPoint)

		// delete keys

		pricePoints.Delete(14)
		pricePoints.Delete(10)
		pricePoints.Delete(13)
		pricePoints.Delete(15)

		pricePoints.SeekToLast()
		pricePoints.SeekToFirst()
		pricePoints.Seek(10)

		// test empty iterator
		iterator = pricePoints.Iterator()
		iterator.Seek(1)
		iterator.Close()

		// test simple iterator
		iterator = pricePoints.Iterator()
		// fivePoint := &PricePoint{Entries: []BookEntry{BookEntry{Amount: 5}}}
		// pricePoints.Set(1, &PricePoint{Entries: []BookEntry{BookEntry{Amount: 1}}})
		// pricePoints.Set(5, fivePoint)
		// iterator.Seek(5)
		// pricePoints.Delete(5)
		// So(iterator.Value(), ShouldEqual, fivePoint)
		// iterator.Next()
		iterator.Seek(1)
		iterator.Close()

		pricePoints.SetProbability(0.1)
		So(pricePoints.probability, ShouldAlmostEqual, 0.1)

		Convey("Should panic at set for key=0", func() {
			So(func() {
				pricePoints := NewPricePoints()
				pricePoints.Set(0, &PricePoint{Entries: []BookEntry{BookEntry{Amount: 1}}})
			}, ShouldPanic)
		})

		Convey("Should panic at delete for key=0", func() {
			So(func() {
				pricePoints := NewPricePoints()
				pricePoints.Delete(0)
			}, ShouldPanic)
		})
	})
}

func TestSkiplists(t *testing.T) {
	Convey("Check skiplist implementation", t, func() {
		skipNode := node{}
		So(skipNode.next(), ShouldBeNil)
		So(maxInt(5, 1), ShouldEqual, 5)
	})
}

func TestCheckCompareFunction(t *testing.T) {
	Convey("Check compare function", t, func() {
		So(cmp_func(10, 3), ShouldBeFalse)
		So(cmp_func(3, 10), ShouldBeTrue)
		So(cmp_func(10, 10), ShouldBeFalse)
	})
}
