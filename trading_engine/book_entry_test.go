package trading_engine_test

import (
	"testing"
	"trading_engine/trading_engine"

	"github.com/francoispqt/gojay"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBookEntryCreation(t *testing.T) {
	Convey("Given an order", t, func() {
		order := trading_engine.NewOrder("TEST_1", 100000000, 12000000000, 1, 1, 1)
		Convey("I should be able to create a new book entry", func() {
			bookEntry := trading_engine.NewBookEntry(order)
			So(bookEntry.Amount, ShouldEqual, 12000000000)
			So(bookEntry.Price, ShouldEqual, 100000000)
			So(bookEntry.Order, ShouldResemble, order)
		})
	})
}

func TestBookEntryLoadFromJson(t *testing.T) {
	Convey("Should be able to load a book entry from json", t, func() {
		var entry trading_engine.BookEntry
		json := []byte(`{
			"price": "1312213.00010201",
			"amount": "8483828.29993942",
			"order": {
				"base": "sym",
				"quote": "tst",
				"id":"TST_1",
				"price": "1312213.00010201",
				"amount": "8483828.29993942",
				"side": 1,
				"type": 1,
				"stop": 1,
				"stop_price": "13132311.00010201",
				"funds": "101000101.33232313"
			}
		}`)
		entry.FromJSON(json)
		So(entry.IsNil(), ShouldBeFalse)
		So(entry.Price, ShouldEqual, 131221300010201)
		So(entry.Amount, ShouldEqual, 848382829993942)
		So(entry.Order.Price, ShouldEqual, 131221300010201)
		So(entry.Order.Amount, ShouldEqual, 848382829993942)
		So(entry.Order.Funds, ShouldEqual, 10100010133232313)
		So(entry.Order.ID, ShouldEqual, "TST_1")
		So(entry.Order.Side, ShouldEqual, 1)
		So(entry.Order.Type, ShouldEqual, 1)
		So(entry.Order.BaseCurrency, ShouldEqual, "sym")
		So(entry.Order.QuoteCurrency, ShouldEqual, "tst")
	})
}

func TestBookEntryConvertToJson(t *testing.T) {
	Convey("Should be able to convert a book entry to json string", t, func() {
		var entry trading_engine.BookEntry
		json := `{"order":{"id":"TST_1","base":"sym","quote":"tst","stop":1,"side":1,"type":1,"event_type":1,"price":"1312213.00010201","amount":"8483828.29993942","stop_price":"13132311.00010201","funds":"101000101.33232313"},"price":"1312213.00010201","amount":"8483828.29993942"}`
		entry.FromJSON([]byte(json))
		bytes, _ := entry.ToJSON()
		So(string(bytes), ShouldEqual, json)
	})
}

func TestBookEntriesConvertFromToJson(t *testing.T) {
	Convey("Should be able to convert a list of book entries to json string", t, func() {
		var entries trading_engine.BookEntries
		json := `[{"order":{"id":"TST_1","base":"sym","quote":"tst","stop":1,"side":1,"type":1,"event_type":1,"price":"1312213.00010201","amount":"8483828.29993942","stop_price":"13132311.00010201","funds":"101000101.33232313"},"price":"1312213.00010201","amount":"8483828.29993942"}]`
		gojay.Unsafe.Unmarshal([]byte(json), &entries)
		So(entries.IsNil(), ShouldEqual, false)
		bytes, _ := gojay.Marshal(entries)
		So(string(bytes), ShouldEqual, json)
	})
}

func TestBookEntriesConvertFromToEmptyJson(t *testing.T) {
	Convey("Should be able to convert a list of book entries to json string", t, func() {
		var entries trading_engine.BookEntries
		json := `[]`
		gojay.Unsafe.Unmarshal([]byte(json), &entries)
		bytes, _ := gojay.Marshal(entries)
		So(string(bytes), ShouldEqual, json)
		So(entries.IsNil(), ShouldEqual, true)
	})
}

func TestBookEntriesConvertFromToBadJson(t *testing.T) {
	Convey("Should not able to convert a list of book entries to json string", t, func() {
		var entries trading_engine.BookEntries
		json := `[{""}]`
		err := gojay.Unsafe.Unmarshal([]byte(json), &entries)
		So(err, ShouldNotBeNil)
		So(entries.IsNil(), ShouldEqual, true)
	})
}
