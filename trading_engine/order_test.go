package trading_engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOrderCreation(t *testing.T) {
	Convey("Should be able to create a new order", t, func() {
		order := NewOrder("TEST_1", 100000000, 12000000000, 1, 2)
		So(order.Amount, ShouldEqual, 12000000000)
		So(order.Price, ShouldEqual, 100000000)
		So(order.Side, ShouldEqual, 1)
		So(order.Type, ShouldEqual, 2)
	})
}

func TestOrderComparisonByPrice(t *testing.T) {
	Convey("Should be able to compare two orders by price", t, func() {
		order1 := NewOrder("TEST_1", 100000000, 12000000000, 1, 2)
		order2 := NewOrder("TEST_2", 110000000, 12000000000, 1, 2)
		So(order1.LessThan(order2), ShouldBeTrue)
	})
}

func TestOrderLoadFromJson(t *testing.T) {
	Convey("Should be able to load an order from json", t, func() {
		var order Order
		json := []byte(`{
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
		}`)
		order.FromJSON(json)
		So(order.IsNil(), ShouldBeFalse)
		So(order.Price, ShouldEqual, 131221300010201)
		So(order.Amount, ShouldEqual, 848382829993942)
		So(order.Funds, ShouldEqual, 10100010133232313)
		So(order.ID, ShouldEqual, "TST_1")
		So(order.Side, ShouldEqual, 1)
		So(order.Type, ShouldEqual, 1)
		So(order.BaseCurrency, ShouldEqual, "sym")
		So(order.QuoteCurrency, ShouldEqual, "tst")
	})
}

func TestOrderConvertToJson(t *testing.T) {
	Convey("Should be able to convert an order to json string", t, func() {
		var order Order
		json := `{"id":"TST_1","base":"sym","quote":"tst","stop":1,"side":1,"type":1,"price":"1312213.00010201","amount":"8483828.29993942","stop_price":"13132311.00010201","funds":"101000101.33232313"}`
		order.FromJSON([]byte(json))
		bytes, _ := order.ToJSON()
		So(string(bytes), ShouldEqual, json)
	})
}
